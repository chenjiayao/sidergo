package datatype

import (
	"errors"
	"strconv"

	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/lib/sortedset"
	"github.com/chenjiayao/goredistraning/redis"
	"github.com/chenjiayao/goredistraning/redis/resp"
)

func init() {

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
