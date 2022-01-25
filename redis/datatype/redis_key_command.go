package datatype

import (
	"strconv"
	"time"

	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/redis"
	"github.com/chenjiayao/goredistraning/redis/resp"
	"github.com/chenjiayao/goredistraning/redis/validate"
)

const (
	UnlimitTTL = int64(-1)
)

func init() {
	redis.RegisterExecCommand(redis.Ttl, ExecTTL, validate.ValidateTtl)
	redis.RegisterExecCommand(redis.Expire, ExecExpire, validate.ValidateExpire)
}

func ExecExpire(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	expire(conn, db, args)
	return resp.MakeNumberResponse(1)
}

// ttl = -2  key 不存在
// ttl = -1 永久有效
func ExecTTL(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	return resp.MakeNumberResponse(ttl(db, args))
}

func ttl(db *redis.RedisDB, args [][]byte) int64 {
	key := string(args[0])

	_, exist := db.Dataset.Get(key)
	if !exist {
		return -2
	}

	//key 存在，但是 TtlMap 中不存在，那么说明key没有设置过期时间
	res, ok := db.TtlMap.Get(key)
	if !ok {
		return -1
	}
	expiredTimestamp, _ := res.(int64)
	now := time.Now().UnixNano() / 1e6
	ttl := (expiredTimestamp - now) / 1000
	if ttl < 0 {
		return -2
	}
	return int64(ttl)
}

//设置key 的 ttl
/*
	保存到 TtlMap 中的是过期的时间
	ttl : 毫秒，以字符串形式传递
*/
func expire(conn conn.Conn, db *redis.RedisDB, args [][]byte) {
	key := string(args[0])
	ttls := string(args[1])

	ttl, _ := strconv.Atoi(ttls)

	if int64(ttl) == UnlimitTTL {
		return
	}

	currentTimestamp := time.Now().UnixNano() / 1e6
	expiredTimestamp := currentTimestamp + int64(ttl)
	db.TtlMap.Put(string(key), expiredTimestamp)
}
