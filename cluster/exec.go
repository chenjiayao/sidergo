package cluster

import (
	"strings"

	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/redis"
	"github.com/chenjiayao/sidergo/redis/resp"
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

	cmdName := strings.ToLower(string(args[0]))
	ipPortPair := cluster.PickNode(cmdName)
	if cluster.Self.IsSelf(ipPortPair) {
		// cluster.Self.RedisServer.
		//TODO这里应该转发给 redisServer 执行命令
		return nil
	} else {
		return nil
	}
}

// mget key1 key2 key3
func ExecMget(cluster *Cluster, args [][]byte) response.Response {
	argsWithoutCmdName := args[1:]
	resps := make([]response.Response, len(argsWithoutCmdName)/2)
	for i := 0; i < len(argsWithoutCmdName); i += 2 {
		getCommand := [][]byte{
			[]byte("get"),
			argsWithoutCmdName[i],
			argsWithoutCmdName[i+1],
		}
		resps = append(resps, defaultExec(cluster, getCommand))
	}
	return resp.MakeArrayResponse(resps)
}

func ExecMset(cluster *Cluster, args [][]byte) response.Response {
	return nil
}
