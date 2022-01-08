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

/*
if cmdName == "auth" {
			res = redisServer.auth(redisClient, args)
			err := redisServer.sendResponse(redisClient, res)
			if err == io.EOF {
				break
			}
			continue
		}

		if !redisServer.isAuthenticated(redisClient) {
			res := resp.MakeErrorResponse("NOAUTH Authentication required")
			err := redisServer.sendResponse(redisClient, res)
			if err == io.EOF {
				break
			}
			continue
		}

		if cmdName == "multi" {
			if redisClient.IsInMultiState() {
				res = resp.MakeErrorResponse("ERR MULTI calls can not be nested")
			} else {
				redisClient.SetMultiState(1)
				res = resp.OKSimpleResponse
			}

			err := redisServer.sendResponse(redisClient, res)
			if err == io.EOF {
				break
			}
		}

		if redisClient.IsInMultiState() {
			redisClient.PushMultiCmd(cmd)
			res = resp.MakeSimpleResponse("QUEUED")
			err := redisServer.sendResponse(redisClient, res)
			if err == io.EOF {
				break
			}
		}

		//执行 select 命令
		if cmdName == "select" {
			dbStr := string(args[0])
			index, err := strconv.Atoi(dbStr)
			if err != nil {
				redisServer.sendResponse(redisClient, resp.MakeErrorResponse("ERR invalid DB index"))
				if err == io.EOF {
					break
				}
			}

			redisClient.SetSelectedDBIndex(index)

			res = resp.MakeSimpleResponse("OK")
			err = redisServer.sendResponse(redisClient, res)
			redisServer.aofHandler.LogCmd(request.Args)
			if err == io.EOF {
				break
			}
			continue
		}
*/
