package redis

import (
	"fmt"
	"sync"

	"github.com/chenjiayao/goredistraning/config"
	"github.com/chenjiayao/goredistraning/interface/db"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/lib/dict"
	"github.com/chenjiayao/goredistraning/redis/resp"
)

var _ db.DB = &RedisDB{}

//aof 规则：
// 1. 不管有多少个 db，只有一个 appendonly.aof 文件，会记录 select db 命令
// 2. 只会记录写命令，不会记录读命令
type RedisDB struct {
	Dataset *dict.ConcurrentDict
	Index   int                  // 数据库 db 编号
	TtlMap  *dict.ConcurrentDict //保存 key 和过期时间之间的关系

	//一个协程来定时删除过期的key
	//一个chan 来关闭「定时删除过期 key 的协程」
	keyLocks sync.Map
}

func NewDBInstance(index int) *RedisDB {
	rd := &RedisDB{
		Dataset:  dict.NewDict(128),
		Index:    index,
		TtlMap:   dict.NewDict(128),
		keyLocks: sync.Map{},
	}
	return rd
}

func (rd *RedisDB) Exec(cmdName string, args [][]byte) response.Response {
	command, ok := CommandTables[cmdName]
	if !ok {
		return resp.MakeErrorResponse(fmt.Sprintf("ERR unknown command '%s'", cmdName))
	}

	//参数校验
	validate := command.ValidateFunc
	if validate != nil {
		err := validate(args)
		if err != nil {
			return resp.MakeErrorResponse(err.Error())
		}
	}

	//执行命令
	execFunc := command.ExecFunc
	resp := execFunc(rd, args)
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
