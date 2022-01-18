package datatype

import (
	"errors"

	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/lib/list"
	"github.com/chenjiayao/goredistraning/redis"
	"github.com/chenjiayao/goredistraning/redis/resp"
	"github.com/chenjiayao/goredistraning/redis/validate"
)

func init() {
	redis.RegisterExecCommand(redis.Lpop, ExecLPop, validate.ValidateLPop)
	redis.RegisterExecCommand(redis.Lpush, ExecLPush, validate.ValidateLPush)
}

func ExecLPop(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	return nil
}

//左到右的顺序依次插入到表头
// lpush key  a b c
//LRANGE mylist 0 -1  ---> c b a
func ExecLPush(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getListOrInitList(conn, db, args)

	if err != nil {
		return resp.MakeErrorResponse("(error) WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	for _, v := range args[1:] {
		l.InsertHead(string(v))
	}
	return resp.MakeNumberResponse(int64(len(args[1:])))

}

func getList(conn conn.Conn, db *redis.RedisDB, args [][]byte) (*list.List, error) {
	key := string(args[0])

	val, exist := db.Dataset.Get(key)
	if exist {
		return nil, nil
	}
	l, ok := val.(*list.List)
	if !ok {
		//TODO报错不是 list 类型
		return nil, errors.New("not list")
	}
	return l, nil
}

func getListOrInitList(conn conn.Conn, db *redis.RedisDB, args [][]byte) (*list.List, error) {
	l, err := getList(conn, db, args)
	if l == nil && err == nil {
		return list.MakeList(), nil
	}
	return l, err
}
