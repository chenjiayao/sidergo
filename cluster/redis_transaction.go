package cluster

import (
	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/request"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/redis"
	"github.com/chenjiayao/sidergo/redis/redisresponse"
	"github.com/chenjiayao/sidergo/redis/validate"
)

func init() {
	RegisterClusterExecCommand(redis.WATCH, ExecWatch, validate.ValidateWatch)
	RegisterClusterExecCommand(redis.MULTI, ExecMulti, validate.ValidateMulti)
	RegisterClusterExecCommand(redis.EXEC, ExecExec, validate.ValidateExec)
	RegisterClusterExecCommand(redis.DISCARD, ExecDiscard, validate.ValidateDiscard)
}

//集群模式下不支持事务，
//其实也能做一做：事务的 key 都应该在同一个 node 下。
func ExecWatch(cluster *Cluster, conn conn.Conn, re request.Request) response.Response {
	return redisresponse.MakeErrorResponse("not support transaction in cluster mode")
}

func ExecMulti(cluster *Cluster, conn conn.Conn, re request.Request) response.Response {
	return redisresponse.MakeErrorResponse("not support transaction in cluster mode")
}

func ExecExec(cluster *Cluster, conn conn.Conn, re request.Request) response.Response {
	return redisresponse.MakeErrorResponse("not support transaction in cluster mode")
}

func ExecDiscard(cluster *Cluster, conn conn.Conn, re request.Request) response.Response {
	return redisresponse.MakeErrorResponse("not support transaction in cluster mode")
}
