package datatype

import (
	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/redis"
	"github.com/chenjiayao/goredistraning/redis/resp"
	"github.com/chenjiayao/goredistraning/redis/validate"
)

func init() {
	redis.RegisterExecCommand(redis.Multi, nil, ExecMulti, nil, validate.ValidateMultiFun)
}

func ExecMulti(conn conn.Conn, args [][]byte) response.Response {
	conn.SetMultiState(1)
	return resp.OKSimpleResponse
}

func ExecDiscard(conn conn.Conn, args [][]byte) response.Response {
	conn.Discard()
	return resp.OKSimpleResponse
}

// watch 的 key ，如果在事务执行之前被其他 client 修改，那么事务不会被执行。
func ExecWatch(conn conn.Conn, args [][]byte) response.Response {
	return nil
}
