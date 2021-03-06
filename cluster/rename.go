package cluster

import (
	"time"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/request"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/redis"
	"github.com/chenjiayao/sidergo/redis/redisresponse"
	"github.com/chenjiayao/sidergo/redis/validate"
)

func init() {

	RegisterClusterExecCommand(redis.RENAME, ExecRename, validate.ValidateRename)
}

//rename oldkey newkey
func ExecRename(cluster *Cluster, conn conn.Conn, re request.Request) response.Response {
	args := re.GetArgs()
	oldKey := string(args[0])
	newKey := string(args[1])

	oldKeyNode := cluster.HashRing.Hit(oldKey)
	newKeyNode := cluster.HashRing.Hit(newKey)
	if oldKeyNode != newKeyNode {
		return redisresponse.MakeErrorResponse("ERR rename must within one slot in cluster mode")
	}

	if cluster.Self.IsSelf(newKeyNode) {
		return cluster.Self.RedisServer.Exec(conn, re)
	} else {
		client := cluster.PeekIdleClient(newKeyNode)
		return client.SendRequest(re, time.Second)
	}
}
