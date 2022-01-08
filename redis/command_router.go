package redis

import (
	"strings"

	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/response"
)

type ExecDBCommandFunc func(db *RedisDB, args [][]byte) response.Response
type ValidateDBCmdArgsFunc func(args [][]byte) error

type ExecConnCommandFunc func(conn conn.Conn, args [][]byte) response.Response
type ValidateConnCmdArgsFunc func(conn conn.Conn, args [][]byte) error

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

	Multi = "multi"

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
	CmdName        string
	DBCommandFun   ExecDBCommandFunc
	ConnCommandFun ExecConnCommandFunc

	DBValidateFunc   ValidateDBCmdArgsFunc
	ConnValidateFunc ValidateConnCmdArgsFunc
}

func RegisterExecCommand(
	cmdName string,
	dbCommandFun ExecDBCommandFunc,
	connCommandFun ExecConnCommandFunc,
	dbValidateFunc ValidateDBCmdArgsFunc,
	connValidateFunc ValidateConnCmdArgsFunc) {

	cmdName = strings.ToLower(cmdName)
	CommandTables[cmdName] = Command{
		CmdName:          cmdName,
		DBCommandFun:     dbCommandFun,
		ConnCommandFun:   connCommandFun,
		DBValidateFunc:   dbValidateFunc,
		ConnValidateFunc: connValidateFunc,
	}
}
