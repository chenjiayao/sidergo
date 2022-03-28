package datatype

import (
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/chenjiayao/sidergo/helper"
	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/redis"
	"github.com/chenjiayao/sidergo/redis/rediserr"
	"github.com/chenjiayao/sidergo/redis/redisresponse"
	"github.com/chenjiayao/sidergo/redis/validate"
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

	redis.RegisterRedisCommand(redis.Set, ExecSet, validate.ValidateSet)
	redis.RegisterRedisCommand(redis.Mget, ExecMGet, validate.ValidateMGet)
	redis.RegisterRedisCommand(redis.Get, ExecGet, validate.ValidateGet)
	redis.RegisterRedisCommand(redis.Incr, ExecIncr, validate.ValidateIncr)
	redis.RegisterRedisCommand(redis.Decr, ExecDecr, validate.ValidateDecr)
	redis.RegisterRedisCommand(redis.Decrby, ExecDecrBy, validate.ValidateDecrBy)
	redis.RegisterRedisCommand(redis.Incrby, ExecIncrBy, validate.ValidateIncrBy)
	redis.RegisterRedisCommand(redis.Incrbyf, ExecIncrByFloat, validate.ValidateIncreByFloat)
	redis.RegisterRedisCommand(redis.Getset, ExecGetset, validate.ValidateGetSet)
	redis.RegisterRedisCommand(redis.Psetex, ExecPSetEX, validate.ValidatePSetEx)
	redis.RegisterRedisCommand(redis.Setnx, ExecSetNX, validate.ValidateSetNx)
	redis.RegisterRedisCommand(redis.Setex, ExecSetEX, validate.ValidateSetEx)
	redis.RegisterRedisCommand(redis.Mset, ExecMSet, validate.ValidateMSet)
	redis.RegisterRedisCommand(redis.Mget, ExecMGet, validate.ValidateMGet)
	redis.RegisterRedisCommand(redis.Msetnx, ExecMSetNX, validate.ValidateMSetNX)
}

func ExecMSet(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	for i := 0; i < len(args); i += 2 {
		ExecSet(conn, db, [][]byte{
			args[i],
			args[i+1],
		})
	}
	return redisresponse.OKSimpleResponse
}

func ExecMGet(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	multiResponses := make([]response.Response, 0)
	for i := 0; i < len(args); i++ {
		r, err := getAsString(conn, db, args[i])
		if err != nil {
			multiResponses = append(multiResponses, redisresponse.NullMultiResponse)
		}
		if r == "" {
			multiResponses = append(multiResponses, redisresponse.MakeMultiResponse(""))
		} else {
			multiResponses = append(multiResponses, redisresponse.MakeMultiResponse(r))
		}
	}
	return redisresponse.MakeArrayResponse(multiResponses)
}

func ExecMSetNX(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	//  这里加锁之前应该对 args 按照字母顺序排序， 保证每个goroutine 的 keys 加锁顺序一致，不然会导致死锁
	allKeys := make([]string, 0)
	for i := 0; i < len(args); i += 2 {
		allKeys = append(allKeys, string(args[i]))
	}
	sort.Slice(allKeys, func(i, j int) bool { return i < j })

	//对所有的 key 加锁
	for i := 0; i < len(allKeys); i++ {
		key := allKeys[i]
		db.LockKey(key, "1")
		defer db.UnLockKey(key)
	}

	//检查是否有哪个 key 已经存在
	for i := 0; i < len(args); i += 2 {
		_, exist := db.Dataset.Get(string(args[i]))
		if exist { //说明 db 中存在这个 key，不管这个 key 是不是 string 类型
			return redisresponse.MakeNumberResponse(0)
		}
	}

	for i := 0; i < len(args); i++ {
		ExecSet(conn, db, [][]byte{
			args[i],
			args[i+1],
		})
	}
	return redisresponse.MakeNumberResponse(1)
}

func ExecGetset(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])

	db.LockKey(key, "1")
	defer db.UnLockKey(key)

	i, exists := db.Dataset.Get(key)

	if !exists {
		return redisresponse.NullMultiResponse
	}

	res, ok := i.(string)
	if !ok {
		return redisresponse.MakeErrorResponse("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	db.Dataset.PutIfExist(key, string(args[1]))
	return redisresponse.MakeSimpleResponse(res)
}

// key value [EX seconds] [PX milliseconds] [NX|XX]
// NX -- Only set the key if it does not already exist.
// XX -- Only set the key if it already exist.   --->同时覆盖新的 ttl
/*
 SET 在设置操作成功完成时，才返回 OK 。
如果设置了 NX 或者 XX ，但因为条件没达到而造成设置操作未执行，那么命令返回空批量回复（NULL Bulk Reply）
*/
func ExecSet(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	ttl := UnlimitTTL

	ss := helper.BbyteToSString(args)
	exFlagIndex := helper.ContainWithoutCaseSensitive(ss, "EX")
	if exFlagIndex != -1 {
		ttlStr := string(args[exFlagIndex+1])
		ttl, _ = strconv.ParseInt(ttlStr, 10, 64)
		ttl *= 1000
	}

	pxFlagIndex := helper.ContainWithoutCaseSensitive(ss, "PX")
	if pxFlagIndex != -1 {
		ttl, _ = strconv.ParseInt(string(args[pxFlagIndex+1]), 10, 64)
	}

	key := string(args[0])
	value := string(args[1])

	//不存在 key 就 insert
	if helper.ContainWithoutCaseSensitive(ss, "NX") != -1 {
		ok := db.Dataset.PutIfNotExist(key, value)
		if ok {
			expire(db, key, ttl)
			return redisresponse.OKSimpleResponse
		} else {
			return redisresponse.NullMultiResponse
		}
	}

	//存在key就 insert
	if helper.ContainWithoutCaseSensitive(ss, "XX") != -1 {
		ok := db.Dataset.PutIfExist(key, value)
		if ok {
			expire(db, key, ttl)
			return redisresponse.OKSimpleResponse
		} else {
			return redisresponse.NullMultiResponse
		}
	}

	ok := db.Dataset.Put(key, value)
	if ok {
		expire(db, key, ttl)
		return redisresponse.OKSimpleResponse
	}
	return redisresponse.NullMultiResponse
}

// setnx key value ---> set key value nx
// 成功返回 1 ，失败返回 0
func ExecSetNX(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	args = append(args, []byte("nx"))
	resp := ExecSet(conn, db, args)
	_, is := resp.(redisresponse.RedisSimpleResponse)
	if is {
		return redisresponse.MakeNumberResponse(1)
	}
	return redisresponse.MakeNumberResponse(0)
}

// setex key seconds value ---> set key value ex second
func ExecSetEX(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	setArgs := [][]byte{
		args[0],
		args[2],
		[]byte("ex"),
		args[1],
	}
	return ExecSet(conn, db, setArgs)
}

// psetex key milliseconds value --> set key value px milliseconds
func ExecPSetEX(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	setArgs := [][]byte{
		args[0],
		args[2],
		[]byte("px"),
		args[1],
	}
	return ExecSet(conn, db, setArgs)
}

/**
get 执行之前要考虑 redis 的过期策略
	redis 的过期策略分为两种方式
		1. 定期删除：每次间隔一定时间再 ttlDict 中扫描，清除过期的 key
		2. 惰性删除：访问一个 key 之前，判断是否已经过期，如果已经过期那么直接删除，并且返回 null
*/
func ExecGet(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	//key 不存在，或者已经到过期时间了
	if ttl(db, [][]byte{args[0]}) < -1 {
		db.Dataset.Del(string(args[0]))
		db.TtlMap.Del(string(args[0]))
		return redisresponse.NullMultiResponse
	}

	s, err := getAsString(conn, db, args[0])
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	if s == "" {
		return redisresponse.NullMultiResponse
	}
	return redisresponse.MakeSimpleResponse(s)
}

func getAsString(conn conn.Conn, db *redis.RedisDB, key []byte) (string, error) {
	res, ok := db.Dataset.Get(string(key))
	if !ok {
		return "", nil
	}
	commo, ok := res.(string)
	if !ok {
		return "", errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	return commo, nil
}

//如果 incr 的key 不存在，那么 set 为1
func ExecIncr(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	incrByArgs := append(args, []byte("1"))
	return ExecIncrBy(conn, db, incrByArgs)
}

func ExecIncrBy(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	steps := string(args[1])
	step, _ := strconv.ParseInt(steps, 10, 64)

	db.LockKey(string(args[0]), "1")
	defer db.UnLockKey(key)

	//get
	s, err := getAsString(conn, db, args[0])
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	val := ""
	//incr
	if s == "" {
		val = fmt.Sprint(step)
	} else {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return redisresponse.MakeErrorResponse(rediserr.NOT_INTEGER_ERROR.Error()) //试图incr 一个字符串
		}
		v += step
		val = fmt.Sprint(v)
	}

	//set
	db.Dataset.Put(key, val)

	return redisresponse.MakeSimpleResponse(val)
}

func ExecIncrByFloat(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	steps := string(args[1])
	step, _ := strconv.ParseFloat(steps, 64)

	db.LockKey(key, "1")
	defer db.UnLockKey(key)

	//get
	s, err := getAsString(conn, db, args[0])
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	val := ""
	//incr
	if s == "" {
		val = fmt.Sprint(step)
	} else {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return redisresponse.MakeErrorResponse(rediserr.NOT_INTEGER_ERROR.Error()) //试图incr 一个字符串
		}
		v += step
		val = fmt.Sprint(v)
	}

	//set
	db.Dataset.Put(key, val)

	return redisresponse.MakeSimpleResponse(val)
}

func ExecDecr(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	incrByArgs := append(args, []byte("-1"))
	return ExecIncrBy(conn, db, incrByArgs)
}

func ExecDecrBy(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	step := string(args[1])
	step = fmt.Sprintf("-%s", step) // 变成 -
	incrByArgs := [][]byte{
		args[0],
		[]byte(step),
	}
	return ExecIncrBy(conn, db, incrByArgs)
}
