package datatype

import (
	"errors"
	"strconv"

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
	redis.RegisterExecCommand(redis.Llen, ExecLLen, validate.ValidateLLen)
	redis.RegisterExecCommand(redis.Lindex, ExecLIndex, validate.ValidateLIndex)
}

func ExecLIndex(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getListOrInitList(conn, db, args)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}
	i, _ := strconv.Atoi(string(args[1]))
	val := l.GetElementByIndex(i)
	if val == nil {
		return resp.NullMultiResponse
	}

	content, _ := val.(string)
	return resp.MakeSimpleResponse(content)
}

//移除并返回列表 key 的头元素。
func ExecLPop(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getListOrInitList(conn, db, args)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}
	element := l.PopFromHead()
	if element == nil {
		return resp.NullMultiResponse
	}

	s, _ := element.(string)
	return resp.MakeSimpleResponse(s)
}

//左到右的顺序依次插入到表头
// lpush key  a b c
//LRANGE mylist 0 -1  ---> c b a
func ExecLPush(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getListOrInitList(conn, db, args)

	if err != nil {
		return resp.MakeErrorResponse(err.Error())
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
		return nil, errors.New("(error) WRONGTYPE Operation against a key holding the wrong kind of value")
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

func ExecLLen(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getListOrInitList(conn, db, args)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}
	return resp.MakeNumberResponse(int64(l.Len()))
}
