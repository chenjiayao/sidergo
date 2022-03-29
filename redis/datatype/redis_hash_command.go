package datatype

import (
	"errors"
	"strconv"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/response"
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
	redis.RegisterRedisCommand(redis.HINCRBY, ExecIncrBy, validate.ValidateIncrBy)
	redis.RegisterRedisCommand(redis.HKEYS, ExecHkeys, validate.ValidateHkeys)
	redis.RegisterRedisCommand(redis.HLEN, ExecHlen, validate.ValidateHlen)
	redis.RegisterRedisCommand(redis.HMGET, ExecHmget, validate.ValidateHmget)
	redis.RegisterRedisCommand(redis.HMSET, ExecHmget, validate.ValidateHmset)
	redis.RegisterRedisCommand(redis.HSETNX, ExecHsetnx, validate.ValidateHsetnx)
	redis.RegisterRedisCommand(redis.HVALS, ExecHvals, validate.ValidateHvals)

}

func ExecHvals(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])

	v, exist := db.Dataset.Get(key)

	if !exist {
		return redisresponse.MakeArrayResponse(nil)
	}

	kvmap := v.(map[string]string)

	multiResponses := make([]response.Response, len(kvmap))
	index := 0
	for _, v := range kvmap {
		multiResponses[index] = redisresponse.MakeMultiResponse(v)
		index++
	}
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

	_, ok := kvmap[field]
	if ok {
		return redisresponse.MakeNumberResponse(0)
	}

	kvmap[field] = value
	return redisresponse.MakeNumberResponse(1)

}
func ExecHmset(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	key := string(args[0])

	kvmap, err := getOrInitHash(db, key)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	for i := 1; i < len(args[:1]); i += 2 {
		field := string(args[i])
		value := string(args[i+1])
		kvmap[field] = value
	}

	return redisresponse.OKSimpleResponse
}

func ExecHmget(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])

	kvmap, err := getOrInitHash(db, key)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	multiResponses := make([]response.Response, len(args[1:]))

	for index, v := range args[1:] {
		value, exist := kvmap[string(v)]
		if !exist {
			multiResponses[index] = redisresponse.NullMultiResponse
		} else {
			multiResponses[index] = redisresponse.MakeMultiResponse(value)
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

	return redisresponse.MakeNumberResponse(int64(len(kvmap)))
}

func ExecHkeys(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	kvmap, err := getOrInitHash(db, string(args[0]))
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	if kvmap == nil {
		return redisresponse.MakeArrayResponse(nil)
	}

	multiResponses := make([]response.Response, len(kvmap))
	index := 0

	for k := range kvmap {
		multiResponses[index] = redisresponse.MakeMultiResponse(k)
		index++
	}

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

	multiResponses := make([]response.Response, len(kvmap)*2)

	index := 0
	for k, v := range kvmap {
		multiResponses[index] = redisresponse.MakeMultiResponse(k)
		multiResponses[index+1] = redisresponse.MakeMultiResponse(v)
		index += 2
	}

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

	value, exist := kvmap[field]
	if !exist {
		value = "0"
	}
	valueAsNumber, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return redisresponse.MakeErrorResponse("ERR hash value is not an integer")
	}

	valueAsNumber += increment
	kvmap[field] = strconv.FormatInt(valueAsNumber, 10)

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

	_, exist := kvmap[field]
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

	_, exist := kvmap[field]

	kvmap[field] = value

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

	value, exist := kvmap[field]
	if !exist {
		return redisresponse.NullMultiResponse
	}
	return redisresponse.MakeMultiResponse(value)
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

		_, exist := kvmap[field]
		if !exist {
			continue
		}
		delete(kvmap, field)
		deletedCount++
	}

	return redisresponse.MakeNumberResponse(deletedCount)
}

/**

在 dataset 中查找 key，尝试转换成 map，如果失败，说明 key 不是 hash 类型，返回 nil,err
如果可以转换 map，需要在 ttl 中查看是否 key 已经过期，如果过期就删除 key，返回 nil, nil
*/
func getOrInitHash(db *redis.RedisDB, key string) (map[string]string, error) {
	v, exist := db.Dataset.Get(key)

	var kvmap map[string]string
	if !exist {
		kvmap = make(map[string]string)
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

		kvmap, ok := v.(map[string]string)

		if !ok {
			return nil, errors.New(" WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		return kvmap, nil
	}
}
