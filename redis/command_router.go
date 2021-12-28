package redis

import (
	"strings"

	"github.com/chenjiayao/goredistraning/interface/response"
)

type ExecFunc func(db *RedisDB, args [][]byte) response.Response
type ValidateCmdArgs func(args [][]byte) bool

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

const (
	set     = "set"
	setnx   = "setnx"
	setex   = "setex"
	psetex  = "psetex"
	mset    = "mset"
	mget    = "mget"
	msetnx  = "msetnx"
	get     = "get"
	getset  = "getset"
	incr    = "incr"
	incrby  = "incrby"
	incrbyf = "incrbyfloat"
	decr    = "decr"
	decrby  = "decrby"
)
