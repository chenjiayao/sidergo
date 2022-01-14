package redis

import (
	"fmt"
	"sync"

	"github.com/chenjiayao/goredistraning/config"
	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/db"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/lib/dict"
	"github.com/chenjiayao/goredistraning/lib/list"
	"github.com/chenjiayao/goredistraning/redis/resp"
)

var _ db.DB = &RedisDB{}

//aof 规则：
// 1. 不管有多少个 db，只有一个 appendonly.aof 文件，会记录 select db 命令
// 2. 只会记录写命令，不会记录读命令
//保存了一个 watched_keys 字典， 字典的键是这个数据库被监视的键， 而字典的值则是一个链表， 链表中保存了所有监视这个键的客户端。
type RedisDB struct {
	Dataset *dict.ConcurrentDict
	Index   int                  // 数据库 db 编号
	TtlMap  *dict.ConcurrentDict //保存 key 和过期时间之间的关系

	//一个协程来定时删除过期的key
	//一个chan 来关闭「定时删除过期 key 的协程」
	keyLocks sync.Map

	// 保存了一个 watched_keys 字典， 字典的键是这个数据库被监视的键， 而字典的值则是一个链表， 链表中保存了所有监视这个键的客户端。
	WatchedKeys sync.Map
}

func NewDBInstance(index int) *RedisDB {
	rd := &RedisDB{
		Dataset:  dict.NewDict(128),
		Index:    index,
		TtlMap:   dict.NewDict(128),
		keyLocks: sync.Map{},

		WatchedKeys: sync.Map{},
	}
	return rd
}

func (rd *RedisDB) canPushMultiQueues(cmdName string) bool {
	canNotPushMultiQueues := map[string]string{
		Exec:    Exec,
		Watch:   Watch,
		Discard: Discard,
		Multi:   Multi,
	}
	_, can := canNotPushMultiQueues[cmdName]
	return !can
}

func (rd *RedisDB) Exec(conn conn.Conn, cmdName string, args [][]byte) response.Response {

	//参数校验
	command, exist := CommandTables[cmdName]
	if !exist {
		return resp.MakeErrorResponse(fmt.Sprintf("ERR unknown command `%s`, with args beginning with:", cmdName))
	}
	validate := command.ValidateFunc
	err := validate(conn, args)
	if err != nil {
		//在 multi 状态下，如果 cmd 校验失败，那么标记 multi 失败，并且返回 error response
		if conn.IsInMultiState() {
			conn.SetMultiState(int(InMultiStateButHaveError))
		}
		return resp.MakeErrorResponse(err.Error())
	}

	//在事务状态下，有些命令不需要 push 到 queue 中s
	if conn.IsInMultiState() && rd.canPushMultiQueues(cmdName) {
		cmd := append([][]byte{[]byte(cmdName)}, args...)
		conn.PushMultiCmd(cmd)
		return resp.MakeSimpleResponse("QUEUED")
	}

	//执行命令
	CommandFunc := command.CommandFunc

	resp := CommandFunc(conn, rd, args)

	_, is := WriteCommands[cmdName]
	if !is {
		return resp
	}

	key := rd.parseCommandKeyFromArgs(args)
	if key != "" {
		rd.setWatchedKeyClientCASDirty(key)
	}
	return resp
}

// 从命令参数中解析出 key
func (rd *RedisDB) parseCommandKeyFromArgs(args [][]byte) string {
	if len(args) == 0 {
		return ""
	}
	key := string(args[0])
	return key
}

//将有 watch key 的 client 的 dirtyCAS 设置为 true
func (rd *RedisDB) setWatchedKeyClientCASDirty(key string) {

	val, exist := rd.WatchedKeys.Load(key)
	if !exist {
		return
	}

	link := val.(*list.List)
	node := link.First()
	for {
		if node == nil {
			break
		}
		v := node.Element()
		conn := v.(*RedisConn)
		conn.DirtyCAS(true)

		node = node.Next()
	}
}

func (rd *RedisDB) LockKey(key string) {
	// 尝试对一个 key 加锁，利用 sync.map 的并发安全特性
	// 但是这里应该挺慢的。。。后续有时间再优化吧
	//for ---> 自旋锁，减少切换线程
	alreadyLockByOtherGoroutine := false
	_, alreadyLockByOtherGoroutine = rd.keyLocks.LoadOrStore(key, 1)
	for alreadyLockByOtherGoroutine {
		_, alreadyLockByOtherGoroutine = rd.keyLocks.LoadOrStore(key, 1)
	}
}

func (rd *RedisDB) UnLockKey(key string) {
	defer rd.keyLocks.Delete(key)
}

func (rd *RedisDB) AddWatchKey(conn conn.Conn, key string) {

	var link *list.List
	val, exist := rd.WatchedKeys.Load(key)
	if !exist {
		link = list.MakeList()
	} else {
		link = val.(*list.List)
	}
	link.InsertIfNotExist(conn)
	rd.WatchedKeys.Store(key, link)
}

func (rd *RedisDB) RemoveWatchKey(conn conn.Conn, key string) {
	var link *list.List
	val, exist := rd.WatchedKeys.Load(key)
	if !exist {
		link = list.MakeList()
	} else {
		link = val.(*list.List)
	}
	link.Remove(conn)
}

func (rd *RedisDB) RemoveAllWatchKey() {
	rd.WatchedKeys = sync.Map{}
}

////////////////
type RedisDBs struct {
	DBs     []*RedisDB
	DBCount int
}

func NewDBs() *RedisDBs {
	dbCount := config.Config.Databases
	rds := &RedisDBs{
		DBs:     make([]*RedisDB, dbCount),
		DBCount: dbCount,
	}

	for i := 0; i < dbCount; i++ {
		rds.DBs[i] = NewDBInstance(i)
	}
	return rds
}
