package cluster

import (
	"time"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/request"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/redis"
	"github.com/chenjiayao/sidergo/redis/redisrequest"
	"github.com/chenjiayao/sidergo/redis/redisresponse"
	"github.com/chenjiayao/sidergo/redis/validate"
	"github.com/sirupsen/logrus"
)

func init() {
	RegisterClusterExecCommand(redis.MGET, ExecMget, validate.ValidateMGet)
	RegisterClusterExecCommand(redis.MSET, ExecMset, validate.ValidateMSet)
	RegisterClusterExecCommand(redis.MSETNX, ExecMSetNX, validate.ValidateMSetNX)

}

// mget key1 key2 key3
func ExecMget(cluster *Cluster, conn conn.Conn, re request.Request) response.Response {
	keys := re.GetArgs()

	resps := make([]response.Response, len(keys))

	for i := 0; i < len(keys); i++ {
		getCommandRequest := &redisrequest.RedisRequet{
			CmdName: redis.GET,
			Args: [][]byte{
				keys[i],
			},
		}
		resps[i] = defaultExec(cluster, conn, getCommandRequest)
	}
	return redisresponse.MakeArrayResponse(resps)
}

func ExecMset(cluster *Cluster, conn conn.Conn, clientRequest request.Request) response.Response {

	args := clientRequest.GetArgs()

	logrus.Info("exec mset")
	undoRequests := make([]request.Request, len(args)/2)
	commitRequests := make([]request.Request, len(args)/2)

	for i := 0; i < len(args); i += 2 {
		undoRequests[i/2] = &redisrequest.RedisRequet{
			CmdName: redis.DEL,
			Args: [][]byte{
				args[i],
			},
		}
		commitRequests[i/2] = &redisrequest.RedisRequet{
			CmdName: redis.SET,
			Args: [][]byte{
				args[i],
				args[i+1],
			},
		}
	}

	tx := MakeTransaction(conn, cluster, undoRequests, commitRequests)
	tx.begin()

	return redisresponse.OKSimpleResponse
}

//msetnx 的所有 key 都应该在同一个 node 中，如果不是那么不执行
func ExecMSetNX(cluster *Cluster, conn conn.Conn, clientRequest request.Request) response.Response {

	args := clientRequest.GetArgs()

	keys := make([]string, len(args)/2)
	for i := 0; i < len(args)/2; i++ {
		keys[i] = string(args[i])
	}

	hitNodeIPPortPair := cluster.HashRing.Hit(string(args[0]))

	for _, k := range keys {
		ipPortPair := cluster.HashRing.Hit(k)
		if hitNodeIPPortPair != ipPortPair {
			return redisresponse.MakeErrorResponse("ERR msetnx must within one slot in cluster mode")
		}
	}

	if cluster.Self.IsSelf(hitNodeIPPortPair) {
		return cluster.Self.RedisServer.Exec(conn, clientRequest)
	}

	client := cluster.PeekIdleClient(hitNodeIPPortPair)
	return client.SendRequest(clientRequest, time.Second)
}

/**

commit'args
	args[0]= set
	args[1] = key
	args[2] = value


rollback's args

*/
