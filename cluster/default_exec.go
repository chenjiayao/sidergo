package cluster

import (
	"fmt"
	"time"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/request"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/redis"
	"github.com/chenjiayao/sidergo/redis/redisrequest"
	"github.com/chenjiayao/sidergo/redis/redisresponse"
	"github.com/chenjiayao/sidergo/redis/validate"
)

func init() {
	RegisterClusterExecCommand(redis.Get, defaultExec, nil)
	RegisterClusterExecCommand(redis.Getset, defaultExec, nil)
	RegisterClusterExecCommand(redis.Set, defaultExec, nil)
	RegisterClusterExecCommand(redis.Incr, defaultExec, nil)
	RegisterClusterExecCommand(redis.Incrbyf, defaultExec, nil)
	RegisterClusterExecCommand(redis.Incrby, defaultExec, nil)

	RegisterClusterExecCommand(redis.Lpush, defaultExec, nil)
	RegisterClusterExecCommand(redis.Lpushx, defaultExec, nil)
	RegisterClusterExecCommand(redis.Rpush, defaultExec, nil)
	RegisterClusterExecCommand(redis.Rpushx, defaultExec, nil)
	RegisterClusterExecCommand(redis.Lpop, defaultExec, nil)
	RegisterClusterExecCommand(redis.Rpop, defaultExec, nil)
	RegisterClusterExecCommand(redis.Lrem, defaultExec, nil)
	RegisterClusterExecCommand(redis.Llen, defaultExec, nil)
	RegisterClusterExecCommand(redis.Lindex, defaultExec, nil)
	RegisterClusterExecCommand(redis.Lset, defaultExec, nil)
	RegisterClusterExecCommand(redis.Lrange, defaultExec, nil)

	RegisterClusterExecCommand(redis.HSET, defaultExec, nil)
	RegisterClusterExecCommand(redis.HSETNX, defaultExec, nil)
	RegisterClusterExecCommand(redis.HGET, defaultExec, nil)
	RegisterClusterExecCommand(redis.HEXISTS, defaultExec, nil)
	RegisterClusterExecCommand(redis.HDEL, defaultExec, nil)
	RegisterClusterExecCommand(redis.HLEN, defaultExec, nil)
	RegisterClusterExecCommand(redis.HMGET, defaultExec, nil)
	RegisterClusterExecCommand(redis.HMSET, defaultExec, nil)
	RegisterClusterExecCommand(redis.HKEYS, defaultExec, nil)
	RegisterClusterExecCommand(redis.HVALS, defaultExec, nil)
	RegisterClusterExecCommand(redis.HGETALL, defaultExec, nil)
	RegisterClusterExecCommand(redis.HINCRBY, defaultExec, nil)
	RegisterClusterExecCommand(redis.HINCRBYFLOAT, defaultExec, nil)

	RegisterClusterExecCommand(redis.Auth, defaultExec, nil)
	RegisterClusterExecCommand(redis.ZREMRANGEBYRANK, defaultExec, nil)
	RegisterClusterExecCommand(redis.ZREMRANGEBYSCORE, defaultExec, nil)
	RegisterClusterExecCommand(redis.ZREM, defaultExec, nil)
	RegisterClusterExecCommand(redis.ZREVRANGEBYSCORE, defaultExec, nil)
	RegisterClusterExecCommand(redis.ZRANGEBYSCORE, defaultExec, nil)
	RegisterClusterExecCommand(redis.ZREVRANGE, defaultExec, nil)
	RegisterClusterExecCommand(redis.ZRANGE, defaultExec, nil)
	RegisterClusterExecCommand(redis.ZCARD, defaultExec, nil)
	RegisterClusterExecCommand(redis.ZREVRANK, defaultExec, nil)
	RegisterClusterExecCommand(redis.ZCOUNT, defaultExec, nil)
	RegisterClusterExecCommand(redis.ZRANK, defaultExec, nil)
	RegisterClusterExecCommand(redis.ZINCRBY, defaultExec, nil)
	RegisterClusterExecCommand(redis.ZSCORE, defaultExec, nil)
	RegisterClusterExecCommand(redis.ZADD, defaultExec, nil)

	RegisterClusterExecCommand(redis.Sadd, defaultExec, nil)
	RegisterClusterExecCommand(redis.Sismember, defaultExec, nil)
	RegisterClusterExecCommand(redis.Scard, defaultExec, nil)
	RegisterClusterExecCommand(redis.Smembers, defaultExec, nil)

	RegisterClusterExecCommand(redis.Select, ExecSelect, validate.ValidateSelect)

	RegisterClusterExecCommand(redis.Ping, ExecPing, validate.ValidatePing)

}

func ExecPing(cluster *Cluster, conn conn.Conn, re request.Request) response.Response {

	args := re.GetArgs()
	message := "PONG"
	if len(args) > 0 {
		message = string(args[0])
	}
	return redisresponse.MakeMultiResponse(message)
}

/*
大部分命令直接走逻辑：
	1. 根据 hashring 定位到 ip:port
	2. 如果定位到的 node 正是当前节点，直接 cluster.Self.RedisServer
	3. 如果不是当前节点，那么通过 TCP 连接转发到对应的节点

	注意，default 的命令中 key 只有一个
*/
func defaultExec(cluster *Cluster, conn conn.Conn, re request.Request) response.Response {

	args := re.GetArgs()
	key := string(args[0])
	ipPortPair := cluster.HashRing.Hit(key)

	if cluster.Self.IsSelf(ipPortPair) {
		return cluster.Self.RedisServer.Exec(conn, re)
	} else {
		c := cluster.PeekIdleClient(ipPortPair)
		selectRequest := &redisrequest.RedisRequet{
			CmdName: "select",
			Args: [][]byte{
				[]byte(fmt.Sprintf("%d", conn.GetSelectedDBIndex())),
			},
		}

		c.SendRequestWithTimeout(selectRequest, 10*time.Second)

		return c.SendRequestWithTimeout(re, 10*time.Second)
	}
}

func ExecSelect(cluster *Cluster, conn conn.Conn, req request.Request) response.Response {

	return cluster.Self.RedisServer.Exec(conn, req)
}
