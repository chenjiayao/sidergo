package cluster

import (
	"strings"

	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/redis"
)

type ClusterExecCommandFunc func(cluster *Cluster, args [][]byte) response.Response
type ClusterValidateFunc redis.ValidateDBCmdArgsFunc

type ClusterCommand struct {
	CmdName      string
	CommandFunc  ClusterExecCommandFunc
	ValidateFunc ClusterValidateFunc
}

var (
	clusterCommandRouter   = make(map[string]ClusterCommand)
	directValidateCommands = map[string]string{
		redis.Mget: "",
	}
)

func RegisterClusterExecCommand(cmdName string, execFn ClusterExecCommandFunc, validateFn redis.ValidateDBCmdArgsFunc) {
	cmdName = strings.ToLower(cmdName)
	clusterCommandRouter[cmdName] = ClusterCommand{
		CmdName:      cmdName,
		ValidateFunc: ClusterValidateFunc(validateFn),
		CommandFunc:  execFn,
	}
}
