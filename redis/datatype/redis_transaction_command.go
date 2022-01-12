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
	redis.RegisterExecCommand(redis.Exec, ExecExec, validate.ValidateExec)

}

func ExecMulti(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	conn.SetMultiState(1)
	return resp.OKSimpleResponse
}

func ExecDiscard(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	conn.Discard()
	db.RemoveAllWatchKey()
	return resp.OKSimpleResponse
}

func ExecExec(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	db.RemoveAllWatchKey()
	conn.SetMultiState(0)
	if conn.GetDirtyCAS() {
		return resp.NullMultiResponse
	}
	return resp.OKSimpleResponse
}

// watch 的 key ，如果在事务执行之前被其他 client 修改，那么事务不会被执行。
// 不管是否已经执行 multi，watch 之后key 被修改，那么事务就不会被执行
// exec 和 discard 执行之后， watch 的 key 被清空
func ExecWatch(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	watchKeys := helper.BbyteToSString(args)
	for i := 0; i < len(watchKeys); i++ {
		watchKey := watchKeys[i]
		db.AddWatchKey(conn, watchKey)
	}
	return resp.OKSimpleResponse
}

func ExecUnwatch(connn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	db.RemoveAllWatchKey()
	return resp.OKSimpleResponse
}
