package redis

import "github.com/chenjiayao/goredistraning/interface/response"

/////// redis string 支持的命令

func (db *RedisDB) StringGet(key string) response.Response {
	return MakeNumberResponse(1)
}
