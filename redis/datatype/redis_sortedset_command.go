package datatype

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/lib/border"
	"github.com/chenjiayao/goredistraning/lib/sortedset"
	"github.com/chenjiayao/goredistraning/redis"
	"github.com/chenjiayao/goredistraning/redis/resp"
	"github.com/chenjiayao/goredistraning/redis/validate"
)

func init() {
	redis.RegisterExecCommand(redis.ZADD, ExecZadd, validate.ValidateZadd)
	redis.RegisterExecCommand(redis.ZCARD, ExecZcard, validate.ValidateZcard)
	redis.RegisterExecCommand(redis.ZCOUNT, ExecZcount, validate.ValidateZcount)
	redis.RegisterExecCommand(redis.ZRANK, ExecZrank, validate.ValidateZrank)
	redis.RegisterExecCommand(redis.ZREVRANGE, ExecZRevrank, validate.ValidateZrevrank)
	redis.RegisterExecCommand(redis.ZREM, ExecZrem, validate.ValidateZrem)
	redis.RegisterExecCommand(redis.ZSCORE, ExecZscore, validate.ValidateZscore)
	redis.RegisterExecCommand(redis.ZINCRBY, ExecZincrby, validate.ValidateIncrBy)
}

//返回有序集 key 中，成员 member 的 score 值。
func ExecZscore(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	ss, err := getAsSortedSet(db, key)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}
	if ss == nil {
		return resp.NullMultiResponse
	}

	member := string(args[1])
	element, exist := ss.Get(member)
	if !exist {
		return resp.NullMultiResponse
	}
	return resp.MakeMultiResponse(fmt.Sprintf("%f", element.Score))
}

func ExecZrem(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	key := string(args[0])
	ss, err := getAsSortedSet(db, key)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}
	if ss == nil {
		return resp.MakeNumberResponse(0)
	}

	for i := 0; i < len(args[1:]); i++ {
		member := string(args[i])
		ss.Remove(member)
	}

	return nil
}

//获取 member 的排名，按照从小到大的顺序
// 排名以 0 为底，score 最小的成员排名为 0
func ExecZrank(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	ss, err := getAsSortedSet(db, key)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}

	member := string(args[1])
	element, exist := ss.Get(member)
	if !exist {
		return resp.NullMultiResponse
	}
	rank := ss.GetRank(element.Memeber, element.Score)

	//正常情况下，不会返回 -1，因为前面已经做过 exist 判断了
	if rank == -1 {
		return resp.NullMultiResponse
	}
	return resp.MakeNumberResponse(rank)
}

func ExecZRevrank(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	ss, err := getAsSortedSet(db, key)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}

	member := string(args[1])
	element, exist := ss.Get(member)
	if !exist {
		return resp.NullMultiResponse
	}
	rank := ss.GetRank(element.Memeber, element.Score)

	//正常情况下，不会返回 -1，因为前面已经做过 exist 判断了
	if rank == -1 {
		return resp.NullMultiResponse
	}
	return resp.MakeNumberResponse(ss.Len() - rank - 1)
}

func ExecZcard(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])

	ss, err := getAsSortedSet(db, key)
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

	key := string(args[0])
	ss, err := getAsSortedSet(db, key)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}
	if ss == nil {
		return resp.MakeNumberResponse(0)
	}

	minBorder, _ := border.ParserBorder(string(args[1]))
	maxBorder, _ := border.ParserBorder(string(args[2]))

	count := ss.Count(minBorder, maxBorder)
	return resp.MakeNumberResponse(count)
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

	element, exist := ss.Get(memeber)
	if !exist {
		ss.Add(memeber, increment)
		return resp.MakeMultiResponse(incrementValue)
	}

	newScore := element.Score + increment
	ss.Remove(element.Memeber)
	ss.Add(element.Memeber, newScore)
	return resp.MakeMultiResponse(fmt.Sprintf("%f", newScore))

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
	ss, err := getAsSortedSet(db, key)
	if err != nil {
		return nil, err
	}
	if ss == nil {
		ss = sortedset.MakeSortedSet()
	}
	return ss, nil
}

func getAsSortedSet(db *redis.RedisDB, key string) (*sortedset.SortedSet, error) {
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
