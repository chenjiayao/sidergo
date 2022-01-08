package datatype

import (
	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/redis"
	"github.com/chenjiayao/goredistraning/redis/validate"
)

func init() {
	redis.RegisterExecCommand(redis.Multi, nil, ExecMulti, nil, validate.ValidateMultiFun)
}

func ExecMulti(conn conn.Conn, args [][]byte) response.Response {
	return nil
}
