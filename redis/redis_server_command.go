package redis

import (
	"strconv"

	"github.com/chenjiayao/goredistraning/interface/response"
)

func init() {
	registerCommand(expire, ExecExpire, nil)
}

func ExecExpire(db *RedisDB, args [][]byte) response.Response {

	resp := ExecGet(db, [][]byte{args[0]})
	if resp == NullMultiResponse {
		return MakeNumberResponse(0)
	}

	ttls := string(args[1])
	ttl, _ := strconv.ParseInt(ttls, 10, 64)

	db.setKeyTtl(args[0], ttl*1000)
	return MakeNumberResponse(1)
}
