package cluster

import (
	"net"

	"github.com/chenjiayao/sidergo/interface/server"
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

///两个 node 节点是否为同一个节点
func (node *Node) IsSelf(n *Node) bool {
	return node.IPAddress == n.IPAddress
}

type Cluster struct {
	self  *Node   //当前 node 节点
	peers []*Node // 集群其他节点
}

func (cluster *Cluster) Handle(conn net.Conn) {

}

func (cluster *Cluster) Close() error {
	cluster.self.RedisServer.Close()
	return nil
}

func (cluster *Cluster) Log() {
	cluster.self.RedisServer.Log()
}

func MakeCluster() *Cluster {

	return nil
}
