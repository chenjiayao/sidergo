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
	redis.RegisterExecCommand(redis.ZRANGE, ExecZrange, validate.ValidateZrange)
	redis.RegisterExecCommand(redis.ZREVRANGE, ExecZrevrange, validate.ValidateZrevrange)
}

func ExecZrevrange(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	key := string(args[0])
	ss, err := getAsSortedSet(db, key)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}
	if ss == nil {
		return resp.EmptyArrayResponse
	}

	withScores := false
	if len(args) == 4 {
		withScores = true
	}
	startValue := string(args[1])
	start, _ := strconv.ParseInt(startValue, 10, 64)

	stopValue := string(args[2])
	stop, _ := strconv.ParseInt(stopValue, 10, 64)

	//将 start stop 的语义转换成 slice 的用法
	if start > ss.Len() || start > stop {
		return resp.EmptyArrayResponse
	}

	//收缩边界
	if start < ss.Len()*-1 {
		start = 0
	}
	if start > ss.Len()-1 {
		start = ss.Len() - 1
	}

	if stop < ss.Len()*-1 {
		stop = -ss.Len()
	}
	if stop > ss.Len()-1 {
		stop = ss.Len() - 1
	}

	if start < 0 {
		start = ss.Len() + start
	}
	if stop < 0 {
		stop = ss.Len() + stop
	}

	elements := ss.Range(start, stop)
	elen := len(elements)

	var responses []response.Response
	if withScores {
		responses = make([]response.Response, len(elements)*2)

		for i := elen - 1; i >= 0; i-- {
			responses[elen-i] = resp.MakeMultiResponse(elements[i].Memeber)
			responses[elen-i+1] = resp.MakeMultiResponse(fmt.Sprintf("%f", elements[i].Score))
		}
	} else {
		responses = make([]response.Response, len(elements))
		for i := 0; i < len(elements); i++ {
			responses[i] = resp.MakeMultiResponse(elements[i].Memeber)
		}
	}
	return resp.MakeArrayResponse(responses)
}

func ExecZrange(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	key := string(args[0])
	ss, err := getAsSortedSet(db, key)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}
	if ss == nil {
		return resp.EmptyArrayResponse
	}

	withScores := false
	if len(args) == 4 {
		withScores = true
	}
	startValue := string(args[1])
	start, _ := strconv.ParseInt(startValue, 10, 64)

	stopValue := string(args[2])
	stop, _ := strconv.ParseInt(stopValue, 10, 64)

	//将 start stop 的语义转换成 slice 的用法
	if start > ss.Len() || start > stop {
		return resp.EmptyArrayResponse
	}

	//收缩边界
	if start < ss.Len()*-1 {
		start = 0
	}
	if start > ss.Len()-1 {
		start = ss.Len() - 1
	}

	if stop < ss.Len()*-1 {
		stop = -ss.Len()
	}
	if stop > ss.Len()-1 {
		stop = ss.Len() - 1
	}

	if start < 0 {
		start = ss.Len() + start
	}
	if stop < 0 {
		stop = ss.Len() + stop
	}

	elements := ss.Range(start, stop)

	var responses []response.Response
	if withScores {
		responses = make([]response.Response, len(elements)*2)
		for i := 0; i < len(elements); i++ {
			responses[i] = resp.MakeMultiResponse(elements[i].Memeber)
			responses[i+1] = resp.MakeMultiResponse(fmt.Sprintf("%f", elements[i].Score))
		}
	} else {
		responses = make([]response.Response, len(elements))
		for i := 0; i < len(elements); i++ {
			responses[i] = resp.MakeMultiResponse(elements[i].Memeber)
		}
	}

	return resp.MakeArrayResponse(responses)
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
