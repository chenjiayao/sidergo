package datatype

import (
	"errors"
	"strconv"
	"strings"

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
	redis.RegisterExecCommand(redis.Lpushx, ExecLPushx, validate.ValidateLPushx)
	redis.RegisterExecCommand(redis.Ltrim, ExecLtrim, validate.ValidateLTrim)
	redis.RegisterExecCommand(redis.Lrange, ExecLrange, validate.ValidateLrange)
	redis.RegisterExecCommand(redis.Linsert, ExecLinsert, validate.ValidateLInsert)
	redis.RegisterExecCommand(redis.Blpop, ExecBlpop, validate.ValidateBlpop)
}

func ExecBlpop(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	return nil
}

func pushGenericCommand(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

}

func ExecLinsert(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getList(conn, db, args)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}
	if l == nil {
		return resp.MakeNumberResponse(0)
	}

	pos := strings.ToUpper(string(args[1])) // before or after
	pivot := string(args[2])
	value := string(args[3])

	var size int64
	if pos == "BEFORE" {
		size = l.InsertBeforePiovt(pivot, value)
	} else if pos == "AFTER" {
		size = l.InsertAfterPiovt(pivot, value)
	}
	return resp.MakeNumberResponse(size)
}

func ExecLrange(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getList(conn, db, args)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}
	if l == nil {
		return resp.EmptyArrayResponse
	}
	start, _ := strconv.Atoi(string(args[1]))
	stop, _ := strconv.Atoi(string(args[2]))

	elements := l.Range(int64(start), int64(stop))

	simpleResponses := make([]response.Response, len(elements))
	for i := 0; i < len(elements); i++ {
		simpleResponses[i] = resp.MakeSimpleResponse(elements[i].(string))
	}
	return resp.MakeArrayResponse(simpleResponses)
}

/**

1. key 可以不存在
2. start 和 stop 两者之间没有任何约束关系
3. start 和 stop 可以是负数
	1. start +，stop -
	2. start +，stop +
	3. start -， stop+
	4. start -， stop-
4. start 和 stop 都是闭区间：[start, stop]
5. start > stop 那么返回空列表，
*/
func ExecLtrim(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	l, err := getList(conn, db, args)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}
	if l == nil {
		return resp.OKSimpleResponse
	}
	start, _ := strconv.Atoi(string(args[1]))
	stop, _ := strconv.Atoi(string(args[2]))
	l.Trim(int64(start), int64(stop))

	return resp.OKSimpleResponse
}

func ExecLPushx(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getList(conn, db, args)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}
	if l == nil {
		return resp.NullMultiResponse
	}

	return ExecLPush(conn, db, args)
}

func ExecLIndex(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getListOrInitList(conn, db, args)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}
	i, _ := strconv.Atoi(string(args[1]))
	val := l.GetElementByIndex(int64(i))
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
	if !exist {
		return nil, nil
	}
	l, ok := val.(*list.List)
	if !ok {
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
