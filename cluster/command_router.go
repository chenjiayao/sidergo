package cluster

import (
	"strings"

	"github.com/chenjiayao/sidergo/interface/response"
)

type ClusterExecCommandFunc func(cluster *Cluster, args [][]byte) response.Response

var clusterCommandRouter = make(map[string]ClusterExecCommandFunc)

func RegisterClusterExecCommand(cmdName string, fn ClusterExecCommandFunc) {
	clusterCommandRouter[strings.ToLower(cmdName)] = fn
}
