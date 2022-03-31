package datatype

import (
	"strconv"

	"github.com/chenjiayao/sidergo/config"
	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/redis"
	"github.com/chenjiayao/sidergo/redis/redisresponse"
	"github.com/chenjiayao/sidergo/redis/validate"
)

func init() {
	redis.RegisterRedisCommand(redis.AUTH, ExecAuth, validate.ValidateAuth)
	redis.RegisterRedisCommand(redis.SELECT, ExecSelect, validate.ValidateSelect)
	redis.RegisterRedisCommand(redis.PERSIST, ExecPersist, validate.ValidatePersist)
	redis.RegisterRedisCommand(redis.EXIST, ExecExist, validate.ValidateExist)
	redis.RegisterRedisCommand(redis.PING, ExecPing, validate.ValidatePing)
}
func ExecExist(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	_, exist := db.Dataset.Get(key)
	if exist {
		return redisresponse.MakeNumberResponse(1)
	}
	return redisresponse.MakeNumberResponse(0)
}

func ExecPersist(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	key := string(args[0])

	_, exist := db.TtlMap.Get(key)
	if !exist {
		return redisresponse.MakeNumberResponse(0)
	}
	db.TtlMap.Del(key)
	return redisresponse.MakeNumberResponse(1)
}
func ExecAuth(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	password := string(args[0])
	if config.Config.RequirePass != password {
		return redisresponse.MakeErrorResponse("ERR invalid password")
	}
	conn.SetPassword(password)
	return redisresponse.MakeSimpleResponse("ok")
}

func ExecSelect(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	dbIndexStr := string(args[0])
	dbIndex, _ := strconv.Atoi(dbIndexStr)
	conn.SetSelectedDBIndex(dbIndex)
	return redisresponse.OKSimpleResponse
}

func ExecPing(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	message := "PONG"
	if len(args) > 0 {
		message = string(args[0])
	}
	return redisresponse.MakeMultiResponse(message)
}
