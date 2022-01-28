package datatype

import (
	"strconv"

	"github.com/chenjiayao/goredistraning/config"
	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/redis"
	"github.com/chenjiayao/goredistraning/redis/resp"
	"github.com/chenjiayao/goredistraning/redis/validate"
)

func init() {
	redis.RegisterExecCommand(redis.Auth, ExecAuth, validate.ValidateAuth)
	redis.RegisterExecCommand(redis.Select, ExecSelect, validate.ValidateSelect)
	redis.RegisterExecCommand(redis.Persist, ExecPersist, validate.ValidatePersist)
}

func ExecPersist(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	key := string(args[0])

	_, exist := db.TtlMap.Get(key)
	if !exist {
		return resp.MakeNumberResponse(0)
	}
	db.TtlMap.Del(key)
	return resp.MakeNumberResponse(1)
}
func ExecAuth(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	password := string(args[0])
	if config.Config.RequirePass != password {
		return resp.MakeErrorResponse("ERR invalid password")
	}
	conn.SetPassword(password)
	return resp.MakeSimpleResponse("ok")
}

func ExecSelect(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	dbIndexStr := string(args[0])
	dbIndex, _ := strconv.Atoi(dbIndexStr)
	conn.SetSelectedDBIndex(dbIndex)
	return resp.OKSimpleResponse
}
