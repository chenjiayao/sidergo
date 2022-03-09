package cluster

import (
	"strings"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/redis"
	redisRequest "github.com/chenjiayao/sidergo/redis/request"
	"github.com/chenjiayao/sidergo/redis/resp"
)

func init() {
	RegisterClusterExecCommand(redis.Get, defaultExec, nil)
	RegisterClusterExecCommand(redis.Getset, defaultExec, nil)
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
}

/*
大部分命令直接走逻辑：
	1. 根据 hashring 定位到 ip:port
	2. 如果定位到的 node 正是当前节点，直接 cluster.Self.RedisServer
	3. 如果不是当前节点，那么通过 TCP 连接转发到对应的节点

	注意，default 的命令中 key 只有一个
*/
func defaultExec(cluster *Cluster, conn conn.Conn, args [][]byte) response.Response {

	key := strings.ToLower(string(args[1]))
	ipPortPair := cluster.HashRing.Hit(key)
	if cluster.Self.IsSelf(ipPortPair) {
		return cluster.Self.RedisServer.Exec(conn, &redisRequest.RedisRequet{
			Args: args,
		})
	} else {
		return nil
	}
}

// mget key1 key2 key3
func ExecMget(cluster *Cluster, conn conn.Conn, args [][]byte) response.Response {
	argsWithoutCmdName := args[1:]
	resps := make([]response.Response, len(argsWithoutCmdName)/2)
	for i := 0; i < len(argsWithoutCmdName); i += 2 {
		getCommand := [][]byte{
			[]byte("get"),
			argsWithoutCmdName[i],
			argsWithoutCmdName[i+1],
		}
		resps = append(resps, defaultExec(cluster, conn, getCommand))
	}
	return resp.MakeArrayResponse(resps)
}

func ExecMset(cluster *Cluster, args [][]byte) response.Response {
	return nil
}
