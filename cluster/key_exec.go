package cluster

import (
	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/request"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/redis"
	req "github.com/chenjiayao/sidergo/redis/request"
	"github.com/chenjiayao/sidergo/redis/resp"
	"github.com/chenjiayao/sidergo/redis/validate"
)

func init() {
	RegisterClusterExecCommand(redis.Del, ExecDel, validate.ValidateDel)
}

func ExecDel(cluster *Cluster, conn conn.Conn, re request.Request) response.Response {
	keys := re.GetArgs()
	resps := make([]response.Response, len(keys))

	for i := 0; i < len(keys); i++ {
		getCommandRequest := &req.RedisRequet{
			CmdName: redis.Get,
			Args: [][]byte{
				keys[i],
			},
		}

		resps[i] = defaultExec(cluster, conn, getCommandRequest)
	}
	return resp.MakeArrayResponse(resps)
}
