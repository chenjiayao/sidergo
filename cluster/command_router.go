package cluster

import (
	"strings"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/redis"
)

type ClusterExecCommandFunc func(cluster *Cluster, conn conn.Conn, args [][]byte) response.Response
type ClusterExecValidateFunc redis.RedisExecValidateFunc

type ClusterCommand struct {
	CmdName      string
	CommandFunc  ClusterExecCommandFunc
	ValidateFunc ClusterExecValidateFunc
}

var (
	clusterCommandRouter = make(map[string]ClusterCommand)

	//需要直接在当前 node 做 validate 的命令
	directValidateCommands = map[string]string{
		redis.Mget: "",
	}
)

func RegisterClusterExecCommand(cmdName string, execFn ClusterExecCommandFunc, validateFn redis.RedisExecValidateFunc) {
	cmdName = strings.ToLower(cmdName)
	clusterCommandRouter[cmdName] = ClusterCommand{
		CmdName:      cmdName,
		ValidateFunc: ClusterExecValidateFunc(validateFn),
		CommandFunc:  execFn,
	}
}
