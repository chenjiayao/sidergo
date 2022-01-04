package datatype

import (
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/lib/set"
	"github.com/chenjiayao/goredistraning/redis"
	"github.com/chenjiayao/goredistraning/redis/resp"
	"github.com/chenjiayao/goredistraning/redis/validate"
)

/**

SADD
SCARD
SMEMBERS


SDIFF
SDIFFSTORE
SINTER
SINTERSTORE
SISMEMBER
SMOVE
SPOP
SRANDMEMBER
SREM
SUNION
SUNIONSTORE
SSCAN
*/
func init() {
	redis.RegisterCommand(redis.Sadd, ExecSadd, validate.ValidateSadd)
	redis.RegisterCommand(redis.Smembers, ExecSmembers, validate.ValidateSmembers)
	redis.RegisterCommand(redis.Scard, ExecScard, validate.ValidateScard)
	redis.RegisterCommand(redis.Spop, ExecSpop, validate.ValidateSpop)
}

const (
	size = 2 >> 64
)

func ExecSpop(db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	s := getSet(db, key)
	l := s.Len()
	if l == 0 {
		return resp.NullMultiResponse
	}

	return nil
}

//SADD runoobkey redis
func ExecSadd(db *redis.RedisDB, args [][]byte) response.Response {

	setValue := getSet(db, string(args[0]))

	for _, v := range args[1:] {
		setValue.Add(string(v))
	}

	db.Dataset.Put(string(args[0]), setValue)
	return resp.MakeNumberResponse(1)
}

func ExecScard(db *redis.RedisDB, args [][]byte) response.Response {

	key := string(args[0])
	s := getSet(db, key)
	return resp.MakeNumberResponse(int64(s.Len()))
}

func ExecSmembers(db *redis.RedisDB, args [][]byte) response.Response {

	setValue := getSet(db, string(args[0]))
	if setValue.Len() == 0 {
		// TODO 返回空数组
		// return resp.make
		return resp.MakeArrayResponse(nil)
	}

	members := setValue.Members()
	return resp.MakeArrayResponse(members)
}

//如果 key 不存在，会新建一个 set
func getSet(db *redis.RedisDB, key string) *set.Set {
	d, exist := db.Dataset.Get(key)

	var setValue *set.Set
	if !exist {
		setValue = set.MakeSet(size)
	} else {
		setValue = d.(*set.Set)
	}
	return setValue
}
