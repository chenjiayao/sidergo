package redis

import (
	"strings"

	"github.com/chenjiayao/goredistraning/interface/response"
)

type ExecFunc func(db *RedisDB, args [][]byte) response.Response
type ValidateCmdArgs func(args [][]byte) error

var (
	CommandTables = make(map[string]Command)
)

type Command struct {
	CmdName      string
	ExecFunc     ExecFunc
	ValidateFunc ValidateCmdArgs
}

func registerCommand(cmdName string, execFunc ExecFunc, validate ValidateCmdArgs) {
	cmdName = strings.ToLower(cmdName)
	CommandTables[cmdName] = Command{
		CmdName:      cmdName,
		ExecFunc:     execFunc,
		ValidateFunc: validate,
	}
}
