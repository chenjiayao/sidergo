package datatype

import (
	"strconv"
	"time"

	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/redis"
	"github.com/chenjiayao/goredistraning/redis/command"
	"github.com/chenjiayao/goredistraning/redis/resp"
)

const (
	UnlimitTTL = int64(-1)
)

func init() {
	redis.RegisterCommand(command.Expire, ExecExpire, nil)
}

func ExecExpire(db *redis.RedisDB, args [][]byte) response.Response {
	ExecTTL(db, args)
	return resp.MakeNumberResponse(1)
}

// ttl = -2  key 不存在
// ttl = -1 永久有效
func ExecTTL(db *redis.RedisDB, args [][]byte) int64 {

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
	expiredAt, _ := res.(int64)
	now := time.Now().UnixNano() / 1e6
	ttl := (expiredAt - now) / 1000
	return int64(ttl)
}

//设置key 的 ttl
/*
	保存到 TtlMap 中的是过期的时间
	ttl : 毫秒，以字符串形式传递
*/
func SetKeyTTL(db *redis.RedisDB, args [][]byte) {
	key := string(args[0])
	ttls := string(args[1])

	ttl, _ := strconv.Atoi(ttls)

	if int64(ttl) == UnlimitTTL {
		return
	}
	expiredAt := time.Now().UnixNano()/1e6 + int64(ttl)
	db.TtlMap.Put(string(key), expiredAt)
}
