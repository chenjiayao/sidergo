package redis

import (
	"sync"

	"github.com/chenjiayao/goredistraning/config"
	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/db"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/lib/dict"
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

func (rd *RedisDB) Exec(conn conn.Conn, cmdName string, args [][]byte) response.Response {
	//参数校验
	command := CommandTables[cmdName]
	validate := command.ValidateFunc

	if validate != nil {
		err := validate(conn, args)
		if err != nil {
			return resp.MakeErrorResponse(err.Error())
		}
	}

	//执行命令
	CommandFunc := command.CommandFunc
	resp := CommandFunc(conn, rd, args)
	return resp
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
