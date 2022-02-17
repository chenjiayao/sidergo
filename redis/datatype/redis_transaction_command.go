package datatype

import (
	"github.com/chenjiayao/sidergo/helper"
	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/redis"
	"github.com/chenjiayao/sidergo/redis/resp"
	"github.com/chenjiayao/sidergo/redis/validate"
)

func init() {

	redis.RegisterExecCommand(redis.Multi, ExecMulti, validate.ValidateMulti)
	redis.RegisterExecCommand(redis.Discard, ExecDiscard, validate.ValidateDiscard)
	redis.RegisterExecCommand(redis.Watch, ExecWatch, validate.ValidateWatch)
	redis.RegisterExecCommand(redis.Exec, ExecExec, validate.ValidateExec)

}

func ExecMulti(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	conn.SetMultiState(int(redis.InMultiState))
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

	if conn.GetMultiState() == int(redis.InMultiStateButHaveError) {
		return resp.MakeErrorResponse("EXECABORT Transaction discarded because of previous errors.")
	}

	multiCmds := conn.GetMultiCmds()

	responseContent := make([]response.Response, len(multiCmds))

	for index, cmd := range multiCmds {
		cmdResponse := db.Exec(conn, string(cmd[0]), cmd[1:])
		if !cmdResponse.ISOK() {
			return resp.MakeErrorResponse("EXECABORT Transaction discarded because of previous errors.")
		}
		responseContent[index] = cmdResponse
	}
	return resp.MakeArrayResponse(responseContent)
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
