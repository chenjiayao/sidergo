package datatype

import (
	"errors"
	"sort"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/lib/set"
	"github.com/chenjiayao/sidergo/redis"
	"github.com/chenjiayao/sidergo/redis/redisresponse"
	"github.com/chenjiayao/sidergo/redis/validate"
)

/**
SADD
SCARD
SMEMBERS
SISMEMBER


SDIFF
SDIFFSTORE
SINTER
SINTERSTORE
SMOVE
SPOP
SRANDMEMBER
SREM
SUNION
SUNIONSTORE
SSCAN
*/
func init() {
	redis.RegisterRedisCommand(redis.Sdiff, ExecSdiff, validate.ValidateSdiff)
	redis.RegisterRedisCommand(redis.Sismember, ExecSismember, validate.ValidateSismember)
	redis.RegisterRedisCommand(redis.Spop, ExecSpop, validate.ValidateSpop)
	redis.RegisterRedisCommand(redis.Sadd, ExecSadd, validate.ValidateSadd)
	redis.RegisterRedisCommand(redis.Scard, ExecScard, validate.ValidateScard)
	redis.RegisterRedisCommand(redis.Smembers, ExecSmembers, validate.ValidateSmembers)
	redis.RegisterRedisCommand(redis.Smove, ExecSmove, validate.ValidateSmove)

}

const (
	size = 2 >> 64
)

/**
smove source destination member

如果 source 没有 member，则 SMOVE 命令不执行任何操作，仅返回 0 。
如果 destination 已经存在 member，则 SMOVE 命令从 source 集合中删除 member 并返回 1 。

*/
func ExecSmove(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	source := string(args[0])
	destination := string(args[1])
	memeber := string(args[2])

	sourceSet, err := getAsSet(conn, db, source)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	destinationSet, err := getSetOrInitSet(conn, db, destination)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	if sourceSet == nil || !sourceSet.Exist(memeber) {
		return redisresponse.MakeNumberResponse(0)
	}

	allKeys := []string{source, destination}
	//顺序加锁，保证不产生死锁
	sort.Slice(allKeys, func(i, j int) bool { return i < j })
	for i := 0; i < len(allKeys); i++ {
		key := allKeys[i]
		db.LockKey(key, "1")
		defer db.UnLockKey(key)
	}
	sourceSet.Del(memeber)
	destinationSet.Add(memeber)
	return redisresponse.MakeNumberResponse(1)
}

func ExecSdiff(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	return nil
}

func ExecSismember(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	s, err := getAsSet(conn, db, key)

	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	if s == nil {
		return redisresponse.MakeNumberResponse(0)
	}

	exist := s.Exist(key)
	if exist {
		return redisresponse.MakeNumberResponse(1)
	}
	return redisresponse.MakeNumberResponse(0)
}

func ExecSpop(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	s, err := getAsSet(conn, db, key)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	if s == nil {
		return redisresponse.NullMultiResponse
	}

	return nil
}

//SADD runoobkey redis
func ExecSadd(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	s, err := getSetOrInitSet(conn, db, string(args[0]))

	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	for _, v := range args[1:] {
		s.Add(string(v))
	}

	db.Dataset.Put(string(args[0]), s)
	return redisresponse.MakeNumberResponse(1)
}

func ExecScard(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	key := string(args[0])
	s, err := getSetOrInitSet(conn, db, key)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	return redisresponse.MakeNumberResponse(int64(s.Len()))
}

func ExecSmembers(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	setValue, err := getAsSet(conn, db, string(args[0]))
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	if setValue.Len() == 0 {
		return redisresponse.MakeArrayResponse(nil)
	}

	members := setValue.Members()
	multiResponses := make([]response.Response, len(members))
	for i := 0; i < len(members); i++ {
		multiResponses[i] = redisresponse.MakeMultiResponse(string(members[i]))
	}
	return redisresponse.MakeArrayResponse(multiResponses)
}

//如果 key 不存在，会新建一个 set
func getSetOrInitSet(conn conn.Conn, db *redis.RedisDB, key string) (*set.Set, error) {
	s, err := getAsSet(conn, db, key)
	if err != nil {
		return nil, err
	}

	if s == nil {
		s = set.MakeSet(size)
		db.Dataset.Put(key, s)
	}
	return s, nil
}

func getAsSet(conn conn.Conn, db *redis.RedisDB, key string) (*set.Set, error) {
	d, exist := db.Dataset.Get(key)
	if exist {
		setValue, ok := d.(*set.Set)
		if !ok {
			return nil, errors.New(" WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		return setValue, nil
	}
	return nil, nil
}
