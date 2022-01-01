package redis

import (
	"strings"

	"github.com/chenjiayao/goredistraning/interface/response"
)

type ExecFunc func(db *RedisDB, args [][]byte) response.Response
type ValidateCmdArgsFunc func(args [][]byte) error

const (
	//string
	Set     = "set"
	Setnx   = "setnx"
	Setex   = "setex"
	Psetex  = "psetex"
	Mset    = "mset"
	Mget    = "mget"
	Msetnx  = "msetnx"
	Get     = "get"
	Getset  = "getset"
	Incr    = "incr"
	Incrby  = "incrby"
	Incrbyf = "incrbyfloat"
	Decr    = "decr"
	Decrby  = "decrby"
	//list
	Lpush     = "lpush"
	Lpushx    = "lpushx"
	Rpush     = "rpush"
	Rpushx    = "rpushx"
	Lpop      = "lpop"
	Rpop      = "rpop"
	Rpoplpush = "rpoplpush"
	Lrem      = "lrem"
	Llen      = "llen"
	Lindex    = "lindex"
	Lset      = "lset"
	Lrange    = "lrange"

	//common
	Expire = "expire"
)

var (
	CommandTables = make(map[string]Command)
)

type Command struct {
	CmdName      string
	ExecFunc     ExecFunc
	ValidateFunc ValidateCmdArgsFunc
}

func RegisterCommand(cmdName string, execFunc ExecFunc, validate ValidateCmdArgsFunc) {
	cmdName = strings.ToLower(cmdName)
	CommandTables[cmdName] = Command{
		CmdName:      cmdName,
		ExecFunc:     execFunc,
		ValidateFunc: validate,
	}
}
