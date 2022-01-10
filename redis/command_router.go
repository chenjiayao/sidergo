package redis

import (
	"strings"

	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/response"
)

type ExecCommandFunc func(conn conn.Conn, db *RedisDB, args [][]byte) response.Response
type ValidateDBCmdArgsFunc func(conn conn.Conn, args [][]byte) error

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

	//set
	Sadd      = "sadd"
	Smembers  = "smembers"
	Scard     = "scard"
	Spop      = "spop"
	Sismember = "sismember"
	Sdiff     = "sdiff"

	Multi   = "multi"
	Discard = "discard"
	Watch   = "watch"

	Auth   = "auth"
	Select = "select"
)

var (
	ConnCommand = map[string]string{
		Multi:  "",
		Select: "",
		Auth:   "",
	}

	DBCommand = map[string]string{
		Set:       "",
		Setnx:     "",
		Mset:      "",
		Msetnx:    "",
		Get:       "",
		Getset:    "",
		Incr:      "",
		Incrby:    "",
		Incrbyf:   "",
		Decr:      "",
		Decrby:    "",
		Lpush:     "",
		Lpushx:    "",
		Rpush:     "",
		Rpushx:    "",
		Lpop:      "",
		Rpop:      "",
		Rpoplpush: "",
		Lrem:      "",
		Llen:      "",
		Lindex:    "",
		Lset:      "",
		Lrange:    "",
		Sadd:      "",
		Smembers:  "",
		Scard:     "",
		Spop:      "",
		Sismember: "",
		Sdiff:     "",
		Expire:    "",
	}
)

var (
	CommandTables = make(map[string]Command)
)

type Command struct {
	CmdName      string
	CommandFunc  ExecCommandFunc
	ValidateFunc ValidateDBCmdArgsFunc
}

func RegisterExecCommand(cmdName string, commandFunc ExecCommandFunc, validateFunc ValidateDBCmdArgsFunc) {

	cmdName = strings.ToLower(cmdName)
	CommandTables[cmdName] = Command{
		CmdName:      cmdName,
		CommandFunc:  commandFunc,
		ValidateFunc: validateFunc,
	}
}
