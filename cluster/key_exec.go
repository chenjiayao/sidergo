package cluster

import (
	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/request"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/redis"
	"github.com/chenjiayao/sidergo/redis/redisrequest"
	"github.com/chenjiayao/sidergo/redis/redisresponse"
	"github.com/chenjiayao/sidergo/redis/validate"
)

func init() {
	RegisterClusterExecCommand(redis.DEL, ExecDel, validate.ValidateDel)
}

func ExecDel(cluster *Cluster, conn conn.Conn, re request.Request) response.Response {
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
