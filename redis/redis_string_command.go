package redis

import (
	"strconv"

	"github.com/chenjiayao/goredistraning/helper"
	"github.com/chenjiayao/goredistraning/interface/response"
)

// - set
// - setnx
// - setex
// - psetex
// - mset
// - mget
// - msetnx
// - get
// - getset
// - incr
// - incrby
// - incrbyfloat
// - decr
// - decrby

func init() {
	registerCommand(set, ExecSet, ValidateSet)
	registerCommand(get, ExecGet, ValidateGet)
}

// key value [EX seconds] [PX milliseconds] [NX|XX]
// NX -- Only set the key if it does not already exist.
// XX -- Only set the key if it already exist.   --->同时覆盖新的 ttl
func ExecSet(db *RedisDB, args [][]byte) response.Response {

	ttl := UnlimitTTL

	ss := helper.BbyteToSString(args)
	exFlagIndex := helper.ContainWithoutCaseSensitive(ss, "EX")
	if exFlagIndex != -1 {
		ttlStr := string(args[exFlagIndex+1])
		ttli, _ := strconv.ParseInt(ttlStr, 10, 64)
		ttl = ttli * 1000
	}

	pxFlagIndex := helper.ContainWithoutCaseSensitive(ss, "PX")
	if pxFlagIndex != -1 {
		ttlStr := string(args[pxFlagIndex+1])
		ttli, _ := strconv.ParseInt(ttlStr, 10, 64)
		ttl = ttli
	}

	key := string(args[0])
	value := string(args[1])

	//不存在 key 就 insert
	if helper.ContainWithoutCaseSensitive(ss, "NX") != -1 {
		ok := db.dataset.PutIfNotExist(key, value)
		if ok {
			db.setKeyTtl(args[0], int64(ttl))
		} else {
			return NullMultiResponse
		}
	}

	//不存在key就 insert
	if helper.ContainWithoutCaseSensitive(ss, "XX") != -1 {
		ok := db.dataset.PutIfExist(key, value)
		if ok {
			db.setKeyTtl(args[0], ttl)
			return OKSimpleResponse
		} else {
			return NullMultiResponse
		}
	}

	ok := db.dataset.Put(key, value)

	if ok {
		db.setKeyTtl(args[0], ttl)
		return OKSimpleResponse
	}
	return NullMultiResponse
}

// setnx key value ---> set key value nx
func ExecSetNX(db *RedisDB, args [][]byte) response.Response {
	args = append(args, []byte("nx"))
	return ExecSet(db, args)
}

func ExecSetEX(db *RedisDB, args [][]byte) response.Response {
	return MakeSimpleResponse("return exec get")
}

func ExecPSetEX(db *RedisDB, args [][]byte) response.Response {
	return MakeSimpleResponse("return exec get")
}

func ExecMSet(db *RedisDB, args [][]byte) response.Response {
	return MakeSimpleResponse("return exec get")

}
func ExecMGet(db *RedisDB, args [][]byte) response.Response {
	return MakeSimpleResponse("return exec get")
}
func ExecGetSet(db *RedisDB, args [][]byte) response.Response {
	return MakeSimpleResponse("return exec get")
}

/**
get 执行之前要考虑 redis 的过期策略
	redis 的过期策略分为两种方式
		1. 定期删除：每次间隔一定时间再 ttlDict 中扫描，清除过期的 key
		2. 惰性删除：访问一个 key 之前，判断是否已经过期，如果已经过期那么直接删除，并且返回 null
*/
func ExecGet(db *RedisDB, args [][]byte) response.Response {

	//key 不存在，或者已经到过期时间了
	if db.ttl(args[0]) < -1 {
		// TODO 删除 key
		return NullMultiResponse
	}

	s := getAsString(db, args[0])
	if s == "" {
		return NullMultiResponse
	}
	return MakeSimpleResponse(s)
}

func getAsString(db *RedisDB, key []byte) string {
	res, ok := db.dataset.Get(string(key))
	if !ok {
		return ""
	}

	commo, ok := res.(string)
	if !ok {
		return ""
	}
	return commo
}

func ExecIncr(db *RedisDB, args [][]byte) response.Response {
	return MakeSimpleResponse("return exec get")

}
func ExecIncrBy(db *RedisDB, args [][]byte) response.Response {
	return MakeSimpleResponse("return exec get")

}
func ExecIncrByFloat(db *RedisDB, args [][]byte) response.Response {
	return MakeSimpleResponse("return exec get")

}
func ExecDecr(db *RedisDB, args [][]byte) response.Response {
	return MakeSimpleResponse("return exec get")

}
func ExecDecrBy(db *RedisDB, args [][]byte) response.Response {
	return MakeSimpleResponse("return exec get")

}
