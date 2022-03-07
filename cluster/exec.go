package cluster

import (
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/redis"
)

func init() {
	RegisterClusterExecCommand(redis.Auth, defaultExec, nil)
}

/*
大部分命令直接走逻辑：
	1. 根据 hashring 定位到 ip:port
	2. 如果定位到的 node 正是当前节点，直接 cluster.Self.RedisServer
	3. 如果不是当前节点，那么通过 TCP 连接转发到对应的节点
*/
func defaultExec(cluster *Cluster, args [][]byte) response.Response {

	return nil
}

// mget key1 key2 key3
func ExecMget(cluster *Cluster, args [][]byte) response.Response {
	argsWithoutCmdName := args[1:]
	for i := 0; i < len(argsWithoutCmdName); i += 2 {
	}
	return nil
}

func ExecMset(cluster *Cluster, args [][]byte) response.Response {
	return nil
}
