package redis

import (
	"github.com/chenjiayao/goredistraning/interface/response"
)

/////// redis string 支持的命令

func (db *RedisDB) ExecGet(key string) response.Response {
	val, ok := db.dataset.Get(key)
	if !ok {
		return MakeErrorResponse("nil")
	}
	v, _ := val.(string)
	return MakeSimpleResponse(v)
}
