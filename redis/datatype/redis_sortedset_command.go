package datatype

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/lib/sortedset"
	"github.com/chenjiayao/goredistraning/redis"
	"github.com/chenjiayao/goredistraning/redis/resp"
	"github.com/chenjiayao/goredistraning/redis/validate"
)

func init() {
	redis.RegisterExecCommand(redis.ZADD, ExecZadd, validate.ValidateZadd)
	redis.RegisterExecCommand(redis.ZCARD, ExecZcard, validate.ValidateZcard)
}

func ExecZcard(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])

	ss, err := getSortedSet(db, key)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}
	if ss == nil {
		return resp.MakeNumberResponse(0)
	}
	return resp.MakeNumberResponse(ss.Len())
}

/*
	zcount key min max
	返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
*/
func ExecZcount(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	return nil
}

/*
	zcount key increment member
	为有序集 key 的成员 member 的 score 值加上增量 increment 。
*/
func ExecZincrby(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	key := string(args[0])
	ss, err := getSortedSetOrInit(db, key)

	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}

	incrementValue := string(args[1])
	increment, _ := strconv.ParseFloat(incrementValue, 64)
	memeber := string(args[2])

	el, exist := ss.Get(memeber)
	if !exist {
		ss.Add(memeber, increment)
		return resp.MakeMultiResponse(incrementValue)
	}
	el.Score += increment
	return resp.MakeMultiResponse(fmt.Sprintf("%f", el.Score))
}

/*
	ZRANGE key start stop [WITHSCORES]
	返回有序集 key 中，指定区间内的成员。
*/
func ExecZrange(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	return nil
}

func ExecZadd(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	key := string(args[0])

	ss, err := getSortedSetOrInit(db, key)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}

	for i := 0; i < len(args[1:]); i += 2 {
		scoreValue := string(args[i])
		member := string(args[i+1])
		score, _ := strconv.ParseFloat(scoreValue, 64)
		ss.Add(member, score)
	}
	return resp.MakeNumberResponse(int64(len(args[1:]) / 2))
}

func getSortedSetOrInit(db *redis.RedisDB, key string) (*sortedset.SortedSet, error) {
	ss, err := getSortedSet(db, key)
	if err != nil {
		return nil, err
	}
	if ss == nil {
		ss = sortedset.MakeSortedSet()
	}
	return ss, nil
}

func getSortedSet(db *redis.RedisDB, key string) (*sortedset.SortedSet, error) {
	entity, exist := db.Dataset.Get(key)
	if !exist {
		return nil, nil
	}

	sortedSet, ok := entity.(*sortedset.SortedSet)
	if !ok {
		return nil, errors.New("-WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	return sortedSet, nil
}
