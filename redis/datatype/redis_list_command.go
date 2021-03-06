package datatype

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/lib/list"
	"github.com/chenjiayao/sidergo/redis"
	"github.com/chenjiayao/sidergo/redis/redisresponse"
	"github.com/chenjiayao/sidergo/redis/validate"
)

func init() {

	redis.RegisterRedisCommand(redis.LLEN, ExecLLen, validate.ValidateLLen)
	redis.RegisterRedisCommand(redis.LINDEX, ExecLIndex, validate.ValidateLIndex)
	redis.RegisterRedisCommand(redis.LTRIM, ExecLtrim, validate.ValidateLTrim)

	redis.RegisterRedisCommand(redis.LRANGE, ExecLrange, validate.ValidateLrange)
	redis.RegisterRedisCommand(redis.LINSERT, ExecLinsert, validate.ValidateLInsert)

	redis.RegisterRedisCommand(redis.LSET, ExecLset, validate.ValidateLset)
	redis.RegisterRedisCommand(redis.LREM, ExecLrem, validate.ValidateLrem)

	redis.RegisterRedisCommand(redis.LPUSH, ExecLPush, validate.ValidateLPush)
	redis.RegisterRedisCommand(redis.RPUSH, ExecRPush, validate.ValidateRPush)

	redis.RegisterRedisCommand(redis.LPOP, ExecLPop, validate.ValidateLPop)
	redis.RegisterRedisCommand(redis.RPOP, ExecRpop, validate.ValidateRPop)

	redis.RegisterRedisCommand(redis.BLPOP, ExecBlpop, validate.ValidateBlpop)
	redis.RegisterRedisCommand(redis.BRPOP, ExecBrpop, validate.ValidateBrpop)

	redis.RegisterRedisCommand(redis.RPUSHX, ExecRPushx, validate.ValidateRPushx)
	redis.RegisterRedisCommand(redis.LPUSHX, ExecLPushx, validate.ValidateLPushx)

	redis.RegisterRedisCommand(redis.RPOPLPUSH, ExecRpoplpush, validate.ValidateRPoplpush)

}

func ExecRpoplpush(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	source := string(args[0])
	destination := string(args[1])

	sourceList, err := getList(conn, db, [][]byte{args[0]})
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	if sourceList == nil {
		return redisresponse.NullMultiResponse
	}

	destinationList, err := getListOrInitList(conn, db, [][]byte{args[1]})
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	var allKeys []string
	if source == destination {
		allKeys = []string{
			source,
		}
	} else {
		allKeys = []string{
			source, destination,
		}
	}

	//????????????????????????????????????
	sort.Slice(allKeys, func(i, j int) bool { return i < j })
	for i := 0; i < len(allKeys); i++ {
		key := allKeys[i]
		db.LockKey(key, "1")
		defer db.UnLockKey(key)
	}

	sourceElement := sourceList.PopFromTail()
	if sourceElement == nil { //source ??????
		return redisresponse.NullMultiResponse
	}

	destinationList.InsertHead(sourceElement)

	val := sourceElement.(string)
	return redisresponse.MakeMultiResponse(val)
}

/**
count > 0 : ?????????????????????????????????????????? value ??????????????????????????? count ???
count < 0 : ?????????????????????????????????????????? value ??????????????????????????? count ???????????????
count = 0 : ????????????????????? value ???????????????

*/
func ExecLrem(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	l, err := getList(conn, db, args)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	if l == nil {
		return redisresponse.MakeNumberResponse(0)
	}

	removeCount := 0
	count, _ := strconv.ParseInt(string(args[1]), 10, 64)
	value := string(args[2])

	if count > l.Len() {
		count = l.Len()
	}
	fromTail := false
	if removeCount < 0 {
		fromTail = true
	}

	node := l.HeadNode()
	if fromTail {
		node = l.TailNode()
	}

	for i := int64(0); i < count; i++ {
		if node.Element() == value {
			l.RemoveNode(node)
			removeCount++
		}
	}
	return redisresponse.MakeNumberResponse(int64(removeCount))
}

func ExecLset(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getList(conn, db, args)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	if l == nil {
		return redisresponse.MakeErrorResponse("ERR no such key")
	}

	index, _ := strconv.ParseInt(string(args[1]), 10, 64)
	if index > l.Len()-1 {
		return redisresponse.MakeErrorResponse("ERR index out of range")
	}

	val := string(args[2])

	node := l.GetNodeByIndex(index)
	node.SetElement(val)
	return redisresponse.OKSimpleResponse
}

func ExecRPushx(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getList(conn, db, args)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	if l == nil {
		return redisresponse.NullMultiResponse
	}
	return ExecRPush(conn, db, args)
}

//????????? pop????????????????????????
/**
1. ?????? db ???????????? blockingKeys ???map map[string]*list.List ???key ??? list ??? key???value ????????????????????????????????? conn
	????????????key ?????????????????????????????? blockingKeys[key] ??????

2. ?????? db ??????????????????  readyList?????? key ??? push ???????????????????????? key ????????? blockingKeys ???????????????????????????????????????(readyList)?????? key ????????????????????? go ????????? channel ??????

*/
func ExecBlpop(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	// ?????????????????? blpop == lpop
	if conn.IsInMultiState() {
		return ExecLPop(conn, db, args)
	}

	timeout, _ := strconv.Atoi(string(args[len(args)-1]))

	for _, v := range args[:len(args)-1] {
		l, err := getList(conn, db, [][]byte{v})
		if err != nil {
			return redisresponse.MakeErrorResponse(err.Error())
		}
		if l != nil {
			v := l.PopFromHead()
			if v != nil {
				content := v.(string)
				return redisresponse.MakeMultiResponse(content)
			}
		}
	}

	//keys ????????????
	conn.SetBlockAt(time.Now())
	conn.SetMaxBlockTime(int64(timeout))
	conn.SetBlockingExec(redis.BLPOP, args)

	for _, v := range args[:len(args)-1] {
		db.AddBlockingConn(string(v), conn)
	}
	return nil
}

func ExecBrpop(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	// ?????????????????? blpop == lpop
	if conn.IsInMultiState() {
		return ExecRpop(conn, db, args)
	}

	timeout, _ := strconv.Atoi(string(args[len(args)-1]))

	for _, v := range args[:len(args)-1] {
		l, err := getList(conn, db, [][]byte{v})
		if err != nil {
			return redisresponse.MakeErrorResponse(err.Error())
		}
		if l != nil {
			v := l.PopFromTail()
			if v != nil {
				content := v.(string)
				return redisresponse.MakeMultiResponse(content)
			}
		}
	}

	//keys ????????????
	conn.SetBlockAt(time.Now())
	conn.SetMaxBlockTime(int64(timeout))
	conn.SetBlockingExec(redis.BLPOP, args)

	for _, v := range args[:len(args)-1] {
		db.AddBlockingConn(string(v), conn)
	}
	return nil
}

//???????????????????????????????????????
// lpush key  a b c
func ExecRPush(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getListOrInitList(conn, db, args)

	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	for _, v := range args[1:] {
		s := string(v)
		l.InsertTail(s)
	}

	db.Dataset.PutIfNotExist(string(args[0]), l)
	db.AddReadyKey(args[0])
	return redisresponse.MakeNumberResponse(l.Len())
}

//???????????????????????????????????????
// lpush key  a b c
//LRANGE mylist 0 -1  ---> c b a
func ExecLPush(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getListOrInitList(conn, db, args)

	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	for _, v := range args[1:] {
		s := string(v)
		l.InsertHead(s)
	}

	db.Dataset.PutIfNotExist(string(args[0]), l)
	db.AddReadyKey(args[0])
	return redisresponse.MakeNumberResponse(l.Len())

}

func ExecLinsert(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getList(conn, db, args)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	if l == nil {
		return redisresponse.MakeNumberResponse(0)
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
	return redisresponse.MakeNumberResponse(size)
}

func ExecLrange(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getList(conn, db, args)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	if l == nil {
		return redisresponse.EmptyArrayResponse
	}
	start, _ := strconv.Atoi(string(args[1]))
	stop, _ := strconv.Atoi(string(args[2]))

	elements := l.Range(int64(start), int64(stop))

	multiResponses := make([]response.Response, len(elements))
	for i := 0; i < len(elements); i++ {
		multiResponses[i] = redisresponse.MakeMultiResponse(elements[i].(string))
	}
	return redisresponse.MakeArrayResponse(multiResponses)
}

/**

1. key ???????????????
2. start ??? stop ????????????????????????????????????
3. start ??? stop ???????????????
	1. start +???stop -
	2. start +???stop +
	3. start -??? stop+
	4. start -??? stop-
4. start ??? stop ??????????????????[start, stop]
5. start > stop ????????????????????????
*/
func ExecLtrim(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {

	l, err := getListOrInitList(conn, db, args)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	if l == nil {
		return redisresponse.OKSimpleResponse
	}
	start, _ := strconv.ParseInt(string(args[1]), 10, 64)
	stop, _ := strconv.ParseInt(string(args[2]), 10, 64)
	l.Trim(start, stop)

	return redisresponse.OKSimpleResponse
}

func ExecLPushx(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getListOrInitList(conn, db, args)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	if l == nil {
		return redisresponse.NullMultiResponse
	}
	return ExecLPush(conn, db, args)
}

func ExecLIndex(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getList(conn, db, args)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}

	if l == nil {
		return redisresponse.NullMultiResponse
	}

	i, _ := strconv.ParseInt(string(args[1]), 10, 64)
	val := l.GetElementByIndex(i)
	if val == nil {
		return redisresponse.NullMultiResponse
	}

	content, _ := val.(string)
	return redisresponse.MakeMultiResponse(content)
}

func ExecRpop(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getListOrInitList(conn, db, args)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	element := l.PopFromTail()
	if element == nil {
		return redisresponse.NullMultiResponse
	}

	s, _ := element.(string)
	return redisresponse.MakeMultiResponse(s)
}

//????????????????????? key ???????????????
func ExecLPop(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getListOrInitList(conn, db, args)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	element := l.PopFromHead()
	if element == nil {
		return redisresponse.NullMultiResponse
	}

	s, _ := element.(string)
	return redisresponse.MakeMultiResponse(s)
}

func ExecLLen(conn conn.Conn, db *redis.RedisDB, args [][]byte) response.Response {
	l, err := getList(conn, db, args)
	if err != nil {
		return redisresponse.MakeErrorResponse(err.Error())
	}
	if l == nil {
		return redisresponse.MakeNumberResponse(0)
	}

	return redisresponse.MakeNumberResponse(l.Len())
}

// get or delete ???????????? getList
func getList(conn conn.Conn, db *redis.RedisDB, args [][]byte) (*list.List, error) {
	key := string(args[0])

	val, exist := db.Dataset.Get(key)
	if !exist {
		return nil, nil
	}
	l, ok := val.(*list.List)
	if !ok {
		return nil, errors.New(" WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	life := ttl(db, [][]byte{
		args[0],
	})

	if life == -2 {
		return nil, nil
	}

	return l, nil
}

// add ???????????? getListOrInitList
func getListOrInitList(conn conn.Conn, db *redis.RedisDB, args [][]byte) (*list.List, error) {
	l, err := getList(conn, db, args)
	if l == nil && err == nil {
		key := string(args[0])
		l := list.MakeList()
		db.Dataset.Put(key, l)
		return l, nil
	}
	return l, err
}
