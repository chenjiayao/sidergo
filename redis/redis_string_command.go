package redis

import (
	"fmt"
	"strconv"

	"github.com/chenjiayao/goredistraning/helper"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/redis/resp"
	"github.com/chenjiayao/goredistraning/rediserr"
)

// - set
// - get
// - incr
// - incrby
// - decr
// - decrby
// - incrbyfloat
// - getset
// - psetex
// - setnx
// - setex
// - mset
// - mget
// - msetnx

func init() {
	registerCommand(set, ExecSet, ValidateSet)
	registerCommand(get, ExecGet, ValidateGet)
	registerCommand(incr, ExecIncr, ValidateIncr)
	registerCommand(incrby, ExecIncrBy, ValidateIncrBy)
	registerCommand(decr, ExecDecr, ValidateDecr)
	registerCommand(decrby, ExecDecrBy, ValidateDecrBy)
	registerCommand(incrbyf, ExecIncrByFloat, ValidateIncreByFloat)
	registerCommand(psetex, ExecPSetEX, ValidatePSetEx)
	registerCommand(getset, ExecGetset, ValidateGetSet)
	registerCommand(setnx, ExecSetNX, ValidateSetNx)
	registerCommand(setex, ExecSetEX, ValidateSetEx)
	registerCommand(mget, ExecMGet, ValidateMGet)
	registerCommand(mset, ExecMSet, ValidateMSet)
	registerCommand(msetnx, ExecMSetNX, validateMSetNX)
}

func ExecMSet(db *RedisDB, args [][]byte) response.Response {
	for i := 0; i < len(args); i += 2 {
		ExecSet(db, [][]byte{
			args[i],
			args[i+1],
		})
	}
	return resp.OKSimpleResponse
}

func ExecMGet(db *RedisDB, args [][]byte) response.Response {

	res := make([][]byte, 0)
	for i := 0; i < len(args); i++ {
		r := getAsString(db, args[i])
		if r == "" {
			res = append(res, nil)
		} else {
			res = append(res, []byte(r))
		}
	}
	return resp.MakeMultiResponse(res)
}

func ExecMSetNX(db *RedisDB, args [][]byte) response.Response {

	//给所有的 key 加锁
	for i := 0; i < len(args); i += 2 {
		key := string(args[i])
		db.lockKey(key)
		defer db.unlockKey(key)
	}

	//检查是否有哪个 key 已经存在
	for i := 0; i < len(args); i += 2 {
		s := getAsString(db, args[i])
		if s != "" {
			return resp.MakeNumberResponse(0)
		}
	}

	for i := 0; i < len(args); i++ {
		ExecSet(db, [][]byte{
			args[i],
			args[i+1],
		})
	}
	return resp.MakeNumberResponse(1)
}

func ExecGetset(db *RedisDB, args [][]byte) response.Response {
	key := string(args[0])

	db.lockKey(key)
	defer db.unlockKey(key)

	i, exists := db.dataset.Get(key)

	if !exists {
		return resp.NullMultiResponse
	}

	res, ok := i.(string)
	if !ok {
		return resp.MakeErrorResponse("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	db.dataset.PutIfExist(key, string(args[1]))
	return resp.MakeSimpleResponse(res)
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
			return resp.NullMultiResponse
		}
	}

	//不存在key就 insert
	if helper.ContainWithoutCaseSensitive(ss, "XX") != -1 {
		ok := db.dataset.PutIfExist(key, value)
		if ok {
			db.setKeyTtl(args[0], ttl)
			return resp.OKSimpleResponse
		} else {
			return resp.NullMultiResponse
		}
	}

	ok := db.dataset.Put(key, value)

	if ok {
		db.setKeyTtl(args[0], ttl)
		return resp.OKSimpleResponse
	}
	return resp.NullMultiResponse
}

// setnx key value ---> set key value nx
func ExecSetNX(db *RedisDB, args [][]byte) response.Response {
	args = append(args, []byte("nx"))
	return ExecSet(db, args)
}

// setex key seconds value ---> set key value ex second
func ExecSetEX(db *RedisDB, args [][]byte) response.Response {
	setArgs := [][]byte{
		args[0],
		args[2],
		[]byte("ex"),
		args[1],
	}
	return ExecSet(db, setArgs)
}

// psetex key milliseconds value --> set key value px milliseconds
func ExecPSetEX(db *RedisDB, args [][]byte) response.Response {
	setArgs := [][]byte{
		args[0],
		args[2],
		[]byte("px"),
		args[1],
	}
	return ExecSet(db, setArgs)
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
		return resp.NullMultiResponse
	}

	s := getAsString(db, args[0])
	if s == "" {
		return resp.NullMultiResponse
	}
	return resp.MakeSimpleResponse(s)
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

//如果 incr 的key 不存在，那么 set 为1
func ExecIncr(db *RedisDB, args [][]byte) response.Response {
	incrByArgs := append(args, []byte("1"))
	return ExecIncrBy(db, incrByArgs)
}

func ExecIncrBy(db *RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	steps := string(args[1])
	step, _ := strconv.ParseInt(steps, 10, 64)

	db.lockKey(string(args[0]))
	defer db.unlockKey(key)

	//get
	s := getAsString(db, args[0])

	val := ""
	//incr
	if s == "" {
		val = fmt.Sprint(step)
	} else {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return resp.MakeErrorResponse(rediserr.NOT_INTEGER_ERROR.Error()) //试图incr 一个字符串
		}
		v += step
		val = fmt.Sprint(v)
	}

	//set
	db.dataset.Put(key, val)

	return resp.MakeSimpleResponse(val)
}
func ExecIncrByFloat(db *RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	steps := string(args[1])
	step, _ := strconv.ParseFloat(steps, 64)

	db.lockKey(key)
	defer db.unlockKey(key)

	//get
	s := getAsString(db, args[0])

	val := ""
	//incr
	if s == "" {
		val = fmt.Sprint(step)
	} else {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return resp.MakeErrorResponse(rediserr.NOT_INTEGER_ERROR.Error()) //试图incr 一个字符串
		}
		v += step
		val = fmt.Sprint(v)
	}

	//set
	db.dataset.Put(key, val)

	return resp.MakeSimpleResponse(val)
}

func ExecDecr(db *RedisDB, args [][]byte) response.Response {
	incrByArgs := append(args, []byte("-1"))
	return ExecIncrBy(db, incrByArgs)
}

func ExecDecrBy(db *RedisDB, args [][]byte) response.Response {

	step := string(args[1])
	step = fmt.Sprintf("-%s", step) // 变成 -
	incrByArgs := [][]byte{
		args[0],
		[]byte(step),
	}
	return ExecIncrBy(db, incrByArgs)
}
