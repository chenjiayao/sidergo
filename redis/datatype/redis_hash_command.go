package datatype

import (
	"errors"
	"strconv"

	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/redis"
	"github.com/chenjiayao/goredistraning/redis/resp"
	"github.com/chenjiayao/goredistraning/redis/validate"
)

/*

hash 的数据结构保存为
	map[key]map[field]value

	TODO 这里应该可以不用 *dict.ConcurrentDict 来做并发
	因为 key get 到的 map 已经加锁了，不会有其他协程可以 get 到这个 key 对应的 map
*/
func init() {

	redis.RegisterExecCommand(redis.HSET, ExecHset, validate.ValidateHset)
	redis.RegisterExecCommand(redis.HGET, ExecHget, validate.ValidateHget)
	redis.RegisterExecCommand(redis.HDEL, ExecHdel, validate.ValidateHdel)
	redis.RegisterExecCommand(redis.HEXISTS, ExecHexists, validate.ValidateHexists)
	redis.RegisterExecCommand(redis.HGETALL, ExecHgetall, validate.ValidateHgetall)
	redis.RegisterExecCommand(redis.HINCRBY, ExecIncrBy, validate.ValidateIncrBy)
	redis.RegisterExecCommand(redis.HKEYS, ExecHkeys, validate.ValidateHkeys)
	redis.RegisterExecCommand(redis.HLEN, ExecHlen, validate.ValidateHlen)
	redis.RegisterExecCommand(redis.HMGET, ExecHmget, validate.ValidateHmget)
	redis.RegisterExecCommand(redis.HMSET, ExecHmget, validate.ValidateHmset)
	redis.RegisterExecCommand(redis.HSETNX, ExecHsetnx, validate.ValidateHsetnx)
	redis.RegisterExecCommand(redis.HVALS, ExecHvals, validate.ValidateHvals)

}

func ExecHvals(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])

	v, exist := db.Dataset.Get(key)

	if !exist {
		return resp.MakeArrayResponse(nil)
	}

	kvmap := v.(map[string]string)

	multiResponses := make([]response.Response, len(kvmap))
	index := 0
	for _, v := range kvmap {
		multiResponses[index] = resp.MakeMultiResponse(v)
		index++
	}
	return resp.MakeArrayResponse(multiResponses)
}

func ExecHsetnx(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	field := string(args[1])
	value := string(args[2])

	kvmap, err := getOrInitHash(db, key)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}

	_, ok := kvmap[field]
	if ok {
		return resp.MakeNumberResponse(0)
	}

	kvmap[field] = value
	return resp.MakeNumberResponse(1)

}
func ExecHmset(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	key := string(args[0])

	kvmap, err := getOrInitHash(db, key)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}

	for i := 1; i < len(args[:1]); i += 2 {
		field := string(args[i])
		value := string(args[i+1])
		kvmap[field] = value
	}

	return resp.OKSimpleResponse
}

func ExecHmget(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])

	kvmap, err := getOrInitHash(db, key)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}

	multiResponses := make([]response.Response, len(args[1:]))

	for index, v := range args[1:] {
		field := string(v)
		value, exist := kvmap[field]
		if !exist {
			multiResponses[index] = resp.NullMultiResponse
		} else {
			multiResponses[index] = resp.MakeMultiResponse(value)
		}
	}
	return resp.MakeArrayResponse(multiResponses)
}

func ExecHlen(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])

	v, exist := db.Dataset.Get(key)

	if !exist {
		return resp.MakeNumberResponse(0)
	}

	kvmap := v.(map[string]string)
	return resp.MakeNumberResponse(int64(len(kvmap)))
}

func ExecHkeys(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])

	v, exist := db.Dataset.Get(key)

	if !exist {
		return resp.MakeArrayResponse(nil)
	}
	kvmap := v.(map[string]string)

	multiResponses := make([]response.Response, len(kvmap))
	index := 0

	for k := range kvmap {
		multiResponses[index] = resp.MakeMultiResponse(k)
		index++
	}

	return resp.MakeArrayResponse(multiResponses)
}

func ExecHgetall(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])

	v, exist := db.Dataset.Get(key)
	if !exist {
		return resp.MakeArrayResponse(nil)
	}
	kvmap := v.(map[string]string)

	multiResponses := make([]response.Response, len(kvmap)*2)

	index := 0
	//TODO 注意，这里使用 map 那么可能每次返回的顺序是不一样的
	for k, v := range kvmap {
		multiResponses[index] = resp.MakeMultiResponse(k)
		multiResponses[index+1] = resp.MakeMultiResponse(v)

		index += 2
	}

	return resp.MakeArrayResponse(multiResponses)

}
func ExecHincrby(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	field := string(args[1])
	increment, _ := strconv.ParseInt(string(args[2]), 10, 64)

	kvmap, err := getOrInitHash(db, key)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}

	value, exist := kvmap[field]
	if !exist {
		value = "0"
	}
	valueAsNumber, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return resp.MakeErrorResponse("ERR hash value is not an integer")
	}

	valueAsNumber += increment
	kvmap[field] = strconv.FormatInt(valueAsNumber, 10)

	return resp.MakeNumberResponse(valueAsNumber)
}

func ExecHexists(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	field := string(args[1])

	v, exist := db.Dataset.Get(key)
	if !exist {
		return resp.MakeNumberResponse(0)
	}

	kvmap := v.(map[string]string)

	_, exist = kvmap[field]
	if exist {
		return resp.MakeNumberResponse(1)
	}
	return resp.MakeNumberResponse(0)

}

func ExecHset(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	key := string(args[0])
	field := string(args[1])
	value := string(args[2])

	kvmap, err := getOrInitHash(db, key)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}

	_, exist := kvmap[field]

	kvmap[field] = value

	if exist {
		return resp.MakeNumberResponse(0)
	}
	return resp.MakeNumberResponse(1)
}

func ExecHget(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	field := string(args[1])

	v, exist := db.Dataset.Get(key)
	if !exist {
		return resp.NullMultiResponse
	}

	kvmap := v.(map[string]string)
	value, exist := kvmap[field]
	if !exist {
		return resp.NullMultiResponse
	}
	return resp.MakeMultiResponse(value)
}

func ExecHdel(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	v, exist := db.Dataset.Get(key)
	if !exist {
		return resp.MakeNumberResponse(0)
	}

	kvmap := v.(map[string]string)

	deletedCount := int64(0)
	for _, v := range args {
		field := string(v)

		_, exist = kvmap[field]
		if !exist {
			continue
		}

		delete(kvmap, field)
		deletedCount++
	}

	return resp.MakeNumberResponse(deletedCount)
}

func getOrInitHash(db *redis.RedisDB, key string) (map[string]string, error) {
	v, exist := db.Dataset.Get(key)

	var kvmap map[string]string
	if !exist {
		// kvmap = dict.NewDict(10)
		kvmap = make(map[string]string)
		db.Dataset.Put(key, kvmap)
		return kvmap, nil
	} else {

		kvmap, ok := v.(map[string]string)

		if !ok {
			return nil, errors.New("(error) WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		return kvmap, nil
	}
}
