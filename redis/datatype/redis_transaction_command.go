package datatype

import (
	"github.com/chenjiayao/goredistraning/helper"
	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/redis"
	"github.com/chenjiayao/goredistraning/redis/resp"
	"github.com/chenjiayao/goredistraning/redis/validate"
)

func init() {

	redis.RegisterExecCommand(redis.Multi, ExecMulti, validate.ValidateMulti)
	redis.RegisterExecCommand(redis.Discard, ExecDiscard, validate.ValidateDiscard)
	redis.RegisterExecCommand(redis.Watch, ExecWatch, validate.ValidateWatch)
}

func ExecMulti(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	conn.SetMultiState(1)
	return resp.OKSimpleResponse
}

func ExecDiscard(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	conn.Discard()
	return resp.OKSimpleResponse
}

func ExecExec(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	return resp.NullMultiResponse
}

// watch 的 key ，如果在事务执行之前被其他 client 修改，那么事务不会被执行。
func ExecWatch(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	watchKeys := helper.BbyteToSString(args)
	for i := 0; i < len(watchKeys); i++ {
		watchKey := watchKeys[i]
		db.AddWatchKey(conn, watchKey)
	}
	return resp.OKSimpleResponse
}
