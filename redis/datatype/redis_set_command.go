package datatype

import (
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

}

const (
	size = 2 >> 64
)

func ExecSdiff(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	return nil
}

func ExecSismember(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	s := getSet(conn, db, key)

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
	s := getSet(conn, db, key)
	if s == nil {
		return redisresponse.NullMultiResponse
	}

	return nil
}

//SADD runoobkey redis
func ExecSadd(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	setValue := getSetOrInitSet(conn, db, string(args[0]))

	for _, v := range args[1:] {
		setValue.Add(string(v))
	}

	db.Dataset.Put(string(args[0]), setValue)
	return redisresponse.MakeNumberResponse(1)
}

func ExecScard(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	key := string(args[0])
	s := getSetOrInitSet(conn, db, key)
	return redisresponse.MakeNumberResponse(int64(s.Len()))
}

func ExecSmembers(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	setValue := getSetOrInitSet(conn, db, string(args[0]))
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
func getSetOrInitSet(conn conn.Conn, db *redis.RedisDB, key string) *set.Set {
	s := getSet(conn, db, key)
	if s == nil {
		return set.MakeSet(size)
	}
	return s
}

func getSet(conn conn.Conn, db *redis.RedisDB, key string) *set.Set {
	d, exist := db.Dataset.Get(key)
	if exist {
		setValue := d.(*set.Set)
		return setValue
	}
	return nil
}
