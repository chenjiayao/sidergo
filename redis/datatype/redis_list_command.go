package datatype

import (
	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/redis"
	"github.com/chenjiayao/goredistraning/redis/validate"
)

func init() {
	redis.RegisterExecCommand(redis.Lpop, ExecLPop, validate.ValidateLPop)
}

func ExecLPop(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	return nil
}
