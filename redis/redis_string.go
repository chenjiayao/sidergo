package redis

import "github.com/chenjiayao/goredistraning/interface/response"

func init() {
	registerCommand(set, ExecSet, ValidateSet)
	registerCommand(get, ExecGet, ValidateGet)
}

func ExecSet(db *RedisDB, args [][]byte) response.Response {
	return MakeSimpleResponse("return exec set")
}

func ExecGet(db *RedisDB, args [][]byte) response.Response {
	return MakeSimpleResponse("return exec get")
}
