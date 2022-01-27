package datatype

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/lib/list"
	"github.com/chenjiayao/goredistraning/redis"
	"github.com/chenjiayao/goredistraning/redis/resp"
	"github.com/chenjiayao/goredistraning/redis/validate"
)

func init() {

	redis.RegisterExecCommand(redis.Llen, ExecLLen, validate.ValidateLLen)
	redis.RegisterExecCommand(redis.Lindex, ExecLIndex, validate.ValidateLIndex)
	redis.RegisterExecCommand(redis.Ltrim, ExecLtrim, validate.ValidateLTrim)
	redis.RegisterExecCommand(redis.Lrange, ExecLrange, validate.ValidateLrange)
	redis.RegisterExecCommand(redis.Linsert, ExecLinsert, validate.ValidateLInsert)
	redis.RegisterExecCommand(redis.Lset, ExecLset, validate.ValidateLset)

	redis.RegisterExecCommand(redis.Lpush, ExecLPush, validate.ValidateLPush)
	redis.RegisterExecCommand(redis.Rpush, ExecRPush, validate.ValidateRPush)

	redis.RegisterExecCommand(redis.Lpop, ExecLPop, validate.ValidateLPop)
	redis.RegisterExecCommand(redis.Rpop, ExecRpop, validate.ValidateRPop)

	redis.RegisterExecCommand(redis.Blpop, ExecBlpop, validate.ValidateBlpop)
	redis.RegisterExecCommand(redis.Brpop, ExecBrpop, validate.ValidateBrpop)

	redis.RegisterExecCommand(redis.Rpushx, ExecRPushx, validate.ValidateRPushx)
	redis.RegisterExecCommand(redis.Lpushx, ExecLPushx, validate.ValidateLPushx)

}

func ExecLset(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getList(conn, db, args)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}
	if l == nil {
		return resp.MakeErrorResponse("(error) ERR no such key")
	}

	index, _ := strconv.ParseInt(string(args[1]), 10, 64)
	if index > l.Len()-1 {
		return resp.MakeErrorResponse("(error) ERR index out of range")
	}

	val := string(args[2])

	node := l.GetNodeByIndex(index)
	node.SetElement(val)
	return resp.MakeSimpleResponse("OK")
}

func ExecRPushx(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getList(conn, db, args)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}
	if l == nil {
		return resp.NullMultiResponse
	}
	return ExecRPush(conn, db, args)
}

func ExecRpop(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getListOrInitList(conn, db, args)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}
	element := l.PopFromTail()
	if element == nil {
		return resp.NullMultiResponse
	}

	s, _ := element.(string)
	return resp.MakeSimpleResponse(s)
}

//右到左的顺序依次插入到表尾
// lpush key  a b c
//LRANGE mylist 0 -1  ---> c b a
func ExecRPush(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getListOrInitList(conn, db, args)

	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}

	for _, v := range args[1:] {
		s := string(v)
		l.InsertTail(s)
	}

	db.Dataset.PutIfNotExist(string(args[0]), l)
	db.AddReadyKey(args[0])
	return resp.MakeNumberResponse(int64(len(args[1:])))
}

//先执行 pop，如果没有，阻塞
/**
1. 每个 db 都有一个 blockingKeys 的map map[string]*list.List ，key 为 list 的 key，value 为链表，保存了被阻塞的 conn
	这样一个key 有没有被阻塞可以通过 blockingKeys[key] 判断

2. 每个 db 也会保存一个  readyList，当 key 被 push 之后，会判断这个 key 在不在 blockingKeys 中，如果在那么创建一个链表(readyList)，将 key 放入链表中（在 go 中使用 channel 代替

*/
func ExecBlpop(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	timeout, _ := strconv.Atoi(string(args[len(args)-1]))

	for _, v := range args[:len(args)-1] {
		l, err := getList(conn, db, [][]byte{v})
		if err != nil {
			return resp.MakeErrorResponse(err.Error())
		}
		if l != nil {
			v := l.PopFromHead()
			if v != nil {
				content := v.(string)
				return resp.MakeSimpleResponse(content)
			}
		}
	}

	//keys 都不存在
	conn.SetBlockAt(time.Now())
	conn.SetMaxBlockTime(int64(timeout))
	conn.SetBlockingExec(redis.Blpop, args)

	for _, v := range args[:len(args)-1] {
		db.AddBlockingConn(string(v), conn)
	}
	return nil
}

func ExecBrpop(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	timeout, _ := strconv.Atoi(string(args[len(args)-1]))

	for _, v := range args[:len(args)-1] {
		l, err := getList(conn, db, [][]byte{v})
		if err != nil {
			return resp.MakeErrorResponse(err.Error())
		}
		if l != nil {
			v := l.PopFromTail()
			if v != nil {
				content := v.(string)
				return resp.MakeSimpleResponse(content)
			}
		}
	}

	//keys 都不存在
	conn.SetBlockAt(time.Now())
	conn.SetMaxBlockTime(int64(timeout))
	conn.SetBlockingExec(redis.Blpop, args)

	for _, v := range args[:len(args)-1] {
		db.AddBlockingConn(string(v), conn)
	}
	return nil
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
		s := string(v)
		l.InsertHead(s)
	}

	db.Dataset.PutIfNotExist(string(args[0]), l)
	db.AddReadyKey(args[0])
	return resp.MakeNumberResponse(int64(len(args[1:])))

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
