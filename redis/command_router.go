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
	Ltrim     = "ltrim"
	Rpoplpush = "rpoplpush"
	Lrem      = "lrem"
	Llen      = "llen"
	Lindex    = "lindex"
	Lset      = "lset"
	Linsert   = "linsert"
	Lrange    = "lrange"
	Blpop     = "blpop"
	Brpop     = "brpop"

	//common
	Expire = "expire"
	Del    = "del"
	Rename = "rename"

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
	Exec    = "exec"

	Auth    = "auth"
	Select  = "select"
	Ttl     = "ttl"
	Persist = "Persist"
)

var (
	CommandTables = make(map[string]Command)

	WriteCommands = map[string]string{
		Set:       "",
		Setnx:     "",
		Setex:     "",
		Persist:   "",
		Psetex:    "",
		Expire:    "",
		Del:       "",
		Lpush:     "",
		Lpushx:    "",
		Rpush:     "",
		Rpushx:    "",
		Lpop:      "",
		Rpop:      "",
		Ltrim:     "",
		Rpoplpush: "",
		Lrem:      "",
		Lindex:    "",
		Lset:      "",
		Linsert:   "",
		Blpop:     "",
		Brpop:     "",

		Mset:    "",
		Msetnx:  "",
		Getset:  "",
		Incr:    "",
		Incrby:  "",
		Incrbyf: "",
		Decr:    "",
		Spop:    "",
		Sadd:    "",
		Decrby:  "",
		Rename:  "",
	}
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
