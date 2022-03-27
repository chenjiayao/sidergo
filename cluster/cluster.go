package cluster

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/chenjiayao/sidergo/config"
	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/request"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/interface/server"
	"github.com/chenjiayao/sidergo/lib/hashring"
	"github.com/chenjiayao/sidergo/redis/redisresponse"
	"github.com/sirupsen/logrus"

	"github.com/chenjiayao/sidergo/parser"
	"github.com/chenjiayao/sidergo/redis"
)

/*
1. 在服务启动之后，检查配置中是否启动集群，如果有，那么创建 Cluster 实例
2. 集群模式下，维护一个环形 hash，每个请求的 key 会映射到某一个 cluster node，如果请求到某一个 node 没有 key，那么会将请求转发到 key 对应的 node
*/

var _ server.Server = &Cluster{}

type Node struct {
	IPAddress   string
	RedisServer *redis.RedisServer
}

func MakeNode(ip string) *Node {
	return &Node{
		IPAddress:   ip,
		RedisServer: nil,
	}
}

///两个 node 节点是否为同一个节点
func (node *Node) IsSelf(ipPortPair string) bool {
	return node.IPAddress == ipPortPair
}

type Cluster struct {
	Self            *Node   //当前 node 节点
	Peers           []*Node // 集群其他节点
	HashRing        *hashring.HashRing
	Pool            map[string]*clientPool
	transactionMaps sync.Map
}

func MakeCluster() *Cluster {

	logrus.Info("enable cluster")
	cluster := &Cluster{
		HashRing:        hashring.MakeHashRing(3),
		Peers:           make([]*Node, len(config.Config.Nodes)),
		transactionMaps: sync.Map{},
		Pool:            make(map[string]*clientPool),
	}

	cluster.Self = MakeNode(config.Config.Self)
	cluster.HashRing.AddNode(config.Config.Self)

	cluster.Self.RedisServer = redis.MakeRedisServer()

	for i := 0; i < len(config.Config.Nodes); i++ {
		ipPortPair := config.Config.Nodes[i]

		node := MakeNode(ipPortPair)
		cluster.Peers[i] = node

		cluster.HashRing.AddNode(ipPortPair)

		//给每个 ipport 配置若干个 client
		cluster.Pool[ipPortPair] = MakeClientPool(ipPortPair, 4)
	}

	return cluster
}

func (cluster *Cluster) Exec(conn conn.Conn, request request.Request) response.Response {

	cmdName := request.GetCmdName()
	args := request.GetArgs()

	command, ok := clusterCommandRouter[cmdName]
	if !ok {
		errResp := redisresponse.MakeErrorResponse(fmt.Sprintf("ERR unknown command `%s`, with args beginning with:", cmdName))
		return errResp
	}

	//在集群模式下，某些命令不是直接转发或者在当前 node 执行，而是要重写逻辑，这部分命令就需要做 validate
	_, ok = directValidateCommands[cmdName]
	if ok {
		err := command.ValidateFunc(conn, args)
		if err != nil {
			errResp := redisresponse.MakeErrorResponse(err.Error())
			return errResp
		}
	}

	res := command.CommandFunc(cluster, conn, request)
	return res
}

/**
1. key 在hash ring 获取节点位置
2. cluster 判断是否在当前节点，如果是直接当前节点处理
3. 如果不是，那么 tcp 转发到对应节点
4. 要考虑分布式事务的处理：MGET，forearch 所有 key 到各个 node 处理，
5. 要考虑 MSET 命令的处理，不可以存在一部分 key set 成功，一部分 set 失败
*/
func (cluster *Cluster) Handle(conn net.Conn) {

	redisClient := redis.MakeRedisConn(conn)
	ch := parser.ReadCommand(conn)
	for request := range ch {
		if request.GetErr() != nil {
			if request.GetErr() == io.EOF {
				cluster.closeClient(redisClient)
				return
			}

			errResponse := redisresponse.MakeErrorResponse(request.GetErr().Error())
			err := redisClient.Write(errResponse.ToErrorByte()) //返回执行命令失败，close client
			if err != nil {
				cluster.closeClient(redisClient)
				return
			}
		}
		if request.GetCmdName() == "" {
			continue
		}

		res := cluster.Exec(redisClient, request)
		err := cluster.sendResponse(redisClient, res)
		if err == io.EOF {
			break
		}
	}
}

func (cluster *Cluster) sendResponse(redisClient conn.Conn, res response.Response) error {

	var err error
	if _, ok := res.(redisresponse.RedisErrorResponse); ok {
		err = redisClient.Write(res.ToErrorByte())
	} else {
		err = redisClient.Write(res.ToContentByte())
	}
	if err == io.EOF {
		cluster.closeClient(redisClient)
	}
	return err
}

func (cluster *Cluster) closeClient(client conn.Conn) {
	client.Close()
}

func (cluster *Cluster) Close() error {
	for _, pool := range cluster.Pool {
		pool.destroy()
	}
	return nil
}

func (cluster *Cluster) Log() {
	cluster.Self.RedisServer.Log()
}

/*

1. mset 要么全部失败，要么全部成功，这个要看看分布式事务
2. mget，mset 命令有多个 key，需要多次 hashring.hit
3. 事务命令   --> 要测试下如果是集群情况下的事务，redis 的表现是怎样的
4. 没有参数的命令
5. 共享登陆状态
6. 在发送命令之前，需要先发送 select 命令
*/

func (cluster *Cluster) PeekIdleClient(ipPortPair string) *client {

	pool := cluster.Pool[ipPortPair]
	for {
		for i := 0; i < len(pool.clients); i++ {
			client := pool.clients[i]
			if client.IsIdle() {
				return client
			}
		}
	}
}
