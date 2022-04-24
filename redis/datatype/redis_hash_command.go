package datatype

import (
	"errors"
	"strconv"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/lib/dict"
	"github.com/chenjiayao/sidergo/redis"
	"github.com/chenjiayao/sidergo/redis/redisresponse"
	"github.com/chenjiayao/sidergo/redis/validate"
)

/*

hash 的数据结构保存为
	map[key]map[field]value

	因为 key get 到的 map 已经加锁了，不会有其他协程可以 get 到这个 key 对应的 map
*/
func init() {

	redis.RegisterRedisCommand(redis.HSET, ExecHset, validate.ValidateHset)
	redis.RegisterRedisCommand(redis.HGET, ExecHget, validate.ValidateHget)
	redis.RegisterRedisCommand(redis.HDEL, ExecHdel, validate.ValidateHdel)
	redis.RegisterRedisCommand(redis.HEXISTS, ExecHexists, validate.ValidateHexists)
	redis.RegisterRedisCommand(redis.HGETALL, ExecHgetall, validate.ValidateHgetall)
	redis.RegisterRedisCommand(redis.HINCRBY, ExecHincrby, validate.ValidateHincrby)
	redis.RegisterRedisCommand(redis.HKEYS, ExecHkeys, validate.ValidateHkeys)
	redis.RegisterRedisCommand(redis.HLEN, ExecHlen, validate.ValidateHlen)
	redis.RegisterRedisCommand(redis.HMGET, ExecHmget, validate.ValidateHmget)
	redis.RegisterRedisCommand(redis.HMSET, ExecHmset, validate.ValidateHmset)
	redis.RegisterRedisCommand(redis.HSETNX, ExecHsetnx, validate.ValidateHsetnx)
	redis.RegisterRedisCommand(redis.HVALS, ExecHvals, validate.ValidateHvals)

}

func ExecHvals(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])

	kvmap, err := getOrInitHash(db, key)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	if kvmap == nil {
		return redisresponse.EmptyArrayResponse
	}

	multiResponses := make([]response.Response, kvmap.Len())
	index := 0

	kvmap.Range(func(key string, val interface{}) {
		multiResponses[index] = redisresponse.MakeMultiResponse(val.(string))
		index++
	})

	return redisresponse.MakeArrayResponse(multiResponses)
}

func ExecHsetnx(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	field := string(args[1])
	value := string(args[2])

	kvmap, err := getOrInitHash(db, key)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	if kvmap == nil {
		kvmap = dict.NewDict(6)
	}

	ok := kvmap.PutIfNotExist(field, value)
	if ok {
		return redisresponse.MakeNumberResponse(1)
	} else {
		return redisresponse.MakeNumberResponse(0)
	}

}
func ExecHmset(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	key := string(args[0])

	kvmap, err := getOrInitHash(db, key)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	if kvmap == nil {
		kvmap = dict.NewDict(6)
	}

	for i := 1; i < len(args[1:]); i += 2 {

		field := string(args[i])
		value := string(args[i+1])
		kvmap.Put(field, value)
	}
	db.Dataset.Put(key, kvmap)

	return redisresponse.OKSimpleResponse
}

func ExecHmget(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])

	kvmap, err := getOrInitHash(db, key)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	if kvmap == nil {
		return redisresponse.EmptyArrayResponse
	}

	multiResponses := make([]response.Response, len(args[1:]))

	for index, v := range args[1:] {
		value, exist := kvmap.Get(string(v))
		if !exist {
			multiResponses[index] = redisresponse.NullMultiResponse
		} else {
			multiResponses[index] = redisresponse.MakeMultiResponse(value.(string))
		}
	}
	return redisresponse.MakeArrayResponse(multiResponses)
}

func ExecHlen(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])

	kvmap, err := getOrInitHash(db, key)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	if kvmap == nil {
		return redisresponse.MakeNumberResponse(0)
	}

	return redisresponse.MakeNumberResponse(int64(kvmap.Len()))
}

func ExecHkeys(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	kvmap, err := getOrInitHash(db, string(args[0]))
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	if kvmap == nil {
		return redisresponse.MakeArrayResponse(nil)
	}

	multiResponses := make([]response.Response, int64(kvmap.Len()))
	index := 0

	kvmap.Range(func(key string, val interface{}) {
		multiResponses[index] = redisresponse.MakeMultiResponse(key)
		index++
	})
	return redisresponse.MakeArrayResponse(multiResponses)
}

func ExecHgetall(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])

	kvmap, err := getOrInitHash(db, key)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	if kvmap == nil {
		return redisresponse.MakeArrayResponse(nil)
	}

	multiResponses := make([]response.Response, int64(kvmap.Len())*2)

	index := 0
	kvmap.Range(func(key string, val interface{}) {
		multiResponses[index] = redisresponse.MakeMultiResponse(key)
		multiResponses[index+1] = redisresponse.MakeMultiResponse(val.(string))
		index += 2
	})

	return redisresponse.MakeArrayResponse(multiResponses)

}
func ExecHincrby(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	field := string(args[1])
	increment, _ := strconv.ParseInt(string(args[2]), 10, 64)

	kvmap, err := getOrInitHash(db, key)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	value, exist := kvmap.Get(field)
	if !exist {
		value = "0"
	}
	valueAsNumber, err := strconv.ParseInt(value.(string), 10, 64)
	if err != nil {
		return redisresponse.MakeErrorResponse("ERR hash value is not an integer")
	}

	valueAsNumber += increment
	kvmap.Put(field, strconv.FormatInt(valueAsNumber, 10))

	return redisresponse.MakeNumberResponse(valueAsNumber)
}

func ExecHexists(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	field := string(args[1])

	kvmap, err := getOrInitHash(db, key)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	if kvmap == nil {
		return redisresponse.MakeNumberResponse(0)
	}

	_, exist := kvmap.Get(field)
	if exist {
		return redisresponse.MakeNumberResponse(1)
	}
	return redisresponse.MakeNumberResponse(0)

}

func ExecHset(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	key := string(args[0])
	field := string(args[1])
	value := string(args[2])

	kvmap, err := getOrInitHash(db, key)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	if kvmap == nil {
		kvmap = dict.NewDict(6)
	}

	_, exist := kvmap.Get(field)

	kvmap.Put(field, value)

	db.Dataset.Put(key, kvmap)

	if exist {
		return redisresponse.MakeNumberResponse(0)
	}

	return redisresponse.MakeNumberResponse(1)
}

func ExecHget(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	field := string(args[1])

	kvmap, err := getOrInitHash(db, key)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	if kvmap == nil {
		return redisresponse.NullMultiResponse
	}

	value, exist := kvmap.Get(field)
	if !exist {
		return redisresponse.NullMultiResponse
	}
	return redisresponse.MakeMultiResponse(value.(string))
}

func ExecHdel(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	kvmap, err := getOrInitHash(db, string(args[0]))
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	if kvmap == nil {
		return redisresponse.MakeNumberResponse(0)
	}

	deletedCount := int64(0)
	for _, v := range args[1:] {
		field := string(v)

		_, exist := kvmap.Get(field)
		if !exist {
			continue
		}
		kvmap.Del(field)
		deletedCount++
	}

	return redisresponse.MakeNumberResponse(deletedCount)
}

/**

在 dataset 中查找 key，尝试转换成 map，如果失败，说明 key 不是 hash 类型，返回 nil,err
如果可以转换 map，需要在 ttl 中查看是否 key 已经过期，如果过期就删除 key，返回 nil, nil
*/
func getOrInitHash(db *redis.RedisDB, key string) (*dict.ConcurrentDict, error) {
	v, exist := db.Dataset.Get(key)

	var kvmap *dict.ConcurrentDict
	if !exist {
		kvmap = dict.NewDict(6)
		db.Dataset.Put(key, kvmap)
		return kvmap, nil
	} else {

		life := ttl(db, [][]byte{
			[]byte(key),
		})

		if life == -2 {
			db.Dataset.Del(key)
			return nil, nil
		}

		kvmap, ok := v.(*dict.ConcurrentDict)

		if !ok {
			return nil, errors.New(" WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		return kvmap, nil
	}
}
