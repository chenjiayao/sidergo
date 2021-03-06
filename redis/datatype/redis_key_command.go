package datatype

import (
	"strconv"
	"time"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/redis"
	"github.com/chenjiayao/sidergo/redis/redisresponse"
	"github.com/chenjiayao/sidergo/redis/validate"
)

const (
	UnlimitTTL = int64(-1)
)

func init() {
	redis.RegisterRedisCommand(redis.TTL, ExecTTL, validate.ValidateTtl)
	redis.RegisterRedisCommand(redis.EXPIRE, ExecExpire, validate.ValidateExpire)
	redis.RegisterRedisCommand(redis.DEL, ExecDel, validate.ValidateDel)
	redis.RegisterRedisCommand(redis.RENAME, ExecRename, validate.ValidateRename)
}

/**
 */
func ExecRename(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	newkey := string(args[1])
	value, exist := db.Dataset.Get(key)
	if !exist {
		return redisresponse.MakeErrorResponse("ERR no such key")
	}

	//不管有没有，删除 newkey 的数据
	db.Dataset.Del(newkey)
	db.TtlMap.Del(newkey)

	//删除 old name key
	db.Dataset.Del(key)

	db.Dataset.Put(newkey, value)

	ttl, exist := db.TtlMap.Get(key)
	if exist {
		db.TtlMap.Put(newkey, ttl)
	}

	return redisresponse.MakeSimpleResponse("OK")
}

func ExecExpire(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	ttls := string(args[1])

	ttl, _ := strconv.Atoi(ttls)

	expire(db, key, int64(ttl*1000))
	return redisresponse.MakeNumberResponse(1)
}

// ttl = -2  key 不存在(已经过期)
// ttl = -1 永久有效
func ExecTTL(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	return redisresponse.MakeNumberResponse(ttl(db, args))
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
	ttl : 毫秒，
*/
func expire(db *redis.RedisDB, key string, ttl int64) {

	if int64(ttl) == UnlimitTTL {
		return
	}
	currentTimestamp := time.Now().UnixNano() / 1e6
	expiredTimestamp := currentTimestamp + ttl
	db.TtlMap.Put(string(key), expiredTimestamp)
}

func ExecDel(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	delCount := 0
	for _, k := range args {
		key := string(k)
		_, exist := db.Dataset.Get(key)
		if exist {
			db.Dataset.Del(key)
			db.TtlMap.Del(key)
			delCount++
		}
	}
	return redisresponse.MakeNumberResponse(int64(delCount))
}
