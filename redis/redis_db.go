package redis

import (
	"fmt"
	"sync"

	"github.com/chenjiayao/sidergo/config"
	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/db"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/interface/server"
	"github.com/chenjiayao/sidergo/lib/dict"
	"github.com/chenjiayao/sidergo/lib/list"
	"github.com/chenjiayao/sidergo/lib/unboundedchan"
	"github.com/chenjiayao/sidergo/redis/redisresponse"
	"github.com/sirupsen/logrus"
)

var _ db.DB = &RedisDB{}

//aof 规则：
// 1. 不管有多少个 db，只有一个 appendonly.aof 文件，会记录 select db 命令
// 2. 只会记录写命令，不会记录读命令
//保存了一个 watched_keys 字典， 字典的键是这个数据库被监视的键， 而字典的值则是一个链表， 链表中保存了所有监视这个键的客户端。
type RedisDB struct {
	Dataset *dict.ConcurrentDict
	Index   int // 数据库 db 编号

	TtlMap *dict.ConcurrentDict //保存 key 和过期时间之间的关系 key ---> unix timestamp

	keyLocks sync.Map //被上锁的key，上锁原因有两个：1. msetnx 需要原子操作，2.集群模式下 prepare 需要先加锁

	WatchedKeys sync.Map // 保存了一个 watched_keys 字典， 字典的键是这个数据库被监视的键， 而字典的值则是一个链表， 链表中保存了所有监视这个键的客户端。

	server server.Server

	BlockingKeys sync.Map                     // key 和被阻塞的「客户端链表」
	ReadyList    *unboundedchan.UnboundedChan // 不再为空的 key 数组

}

func NewDBInstance(server server.Server, index int) *RedisDB {
	rd := &RedisDB{
		Dataset:      dict.NewDict(2),
		Index:        index,
		TtlMap:       dict.NewDict(2),
		keyLocks:     sync.Map{},
		BlockingKeys: sync.Map{},
		WatchedKeys:  sync.Map{},
		server:       server,

		ReadyList: unboundedchan.MakeUnboundedChan(20),
	}

	go rd.handleClientsBlockedOnLists()
	return rd
}

func (rd *RedisDB) CloseDB() {
	logrus.Info(rd.Index, " db closed")
	close(rd.ReadyList.In)
}

func (rd *RedisDB) canPushMultiQueues(cmdName string) bool {
	canNotPushMultiQueues := map[string]string{
		EXEC:    EXEC,
		WATCH:   WATCH,
		DISCARD: DISCARD,
		MULTI:   MULTI,
	}
	_, can := canNotPushMultiQueues[cmdName]
	return !can
}

func (rd *RedisDB) Exec(conn conn.Conn, cmdName string, args [][]byte) response.Response {

	//参数校验
	command, exist := CommandTables[cmdName]
	if !exist {
		return redisresponse.MakeErrorResponse(fmt.Sprintf("ERR unknown command `%s`, with args beginning with:", cmdName))
	}
	validate := command.ValidateFunc
	err := validate(conn, args)
	if err != nil {
		//在 multi 状态下，如果 cmd 校验失败，那么标记 multi 失败，并且返回 error response
		if conn.IsInMultiState() {
			conn.SetMultiState(int(InMultiStateButHaveError))
		}
		return redisresponse.MakeErrorResponse(err.Error())
	}

	//在事务状态下，有些命令不需要 push 到 queue 中s
	if conn.IsInMultiState() && rd.canPushMultiQueues(cmdName) {
		cmd := append([][]byte{[]byte(cmdName)}, args...)
		conn.PushMultiCmd(cmd)
		return redisresponse.MakeMultiResponse("QUEUED")
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
	node := link.HeadNode()
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

func (rd *RedisDB) LockKey(key string, placeholder string) {
	// 尝试对一个 key 加锁，利用 sync.map 的并发安全特性
	// 但是这里应该挺慢的。。。后续有时间再优化吧
	//for ---> 自旋锁，减少切换线程
	alreadyLockByOtherGoroutine := false
	_, alreadyLockByOtherGoroutine = rd.keyLocks.LoadOrStore(key, placeholder)
	for alreadyLockByOtherGoroutine {
		_, alreadyLockByOtherGoroutine = rd.keyLocks.LoadOrStore(key, placeholder)
	}
}

func (rd *RedisDB) UnLockKey(key string) {
	defer rd.keyLocks.Delete(key)
}

////////////////事务相关命令支持//////////////////
func (rd *RedisDB) RemoveWatchKey(conn conn.Conn, key string) {
	var link *list.List
	val, exist := rd.WatchedKeys.Load(key)
	if !exist {
		link = list.MakeList()
	} else {
		link = val.(*list.List)
	}
	link.RemoveNode(conn)
}

func (rd *RedisDB) RemoveAllWatchKey() {
	rd.WatchedKeys = sync.Map{}
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

////////////////事务相关命令支持//////////////////

////////////////list block command 支持//////////////////
func (rd *RedisDB) AddBlockingConn(key string, conn conn.Conn) {

	v, ok := rd.BlockingKeys.Load(key)
	var l *list.List
	if !ok {
		l = list.MakeList()
	} else {
		l = v.(*list.List)
	}
	l.InsertTail(conn)
	rd.BlockingKeys.Store(key, l)
}

func (rd *RedisDB) AddReadyKey(key []byte) {
	k := string(key)
	_, ok := rd.BlockingKeys.Load(k) //没有 conn 因为 key 不存在阻塞
	if !ok {
		return
	}
	//insert
	rd.ReadyList.In <- [][]byte{key}
}

////////////////list block command 支持//////////////////

func (rd *RedisDB) handleClientsBlockedOnLists() {
	for o := range rd.ReadyList.Out {
		if len(o) == 0 {
			continue
		}
		key := string(o[0])
		lv, ok := rd.BlockingKeys.Load(key)
		if !ok {
			continue
		}
		l, _ := lv.(*list.List)
		cv := l.GetElementByIndex(0)
		if cv == nil {
			continue
		}
		conn := cv.(conn.Conn)
		cmdName, args := conn.GetBlockingExec()

		command := CommandTables[cmdName]
		res := command.CommandFunc(conn, rd, args)

		if res == nil {
			continue
		}

		conn.SetBlockingResponse(res) //设置回复
		conn.SetMaxBlockTime(0)
		conn.SetBlockingExec("", nil)
		l.PopFromHead() // conn 不再阻塞
	}
}

////////////////
type RedisDBs struct {
	DBs     []*RedisDB
	DBCount int
}

func NewDBs(server server.Server) *RedisDBs {
	dbCount := config.Config.Databases
	rds := &RedisDBs{
		DBs:     make([]*RedisDB, dbCount),
		DBCount: dbCount,
	}
	for i := 0; i < dbCount; i++ {
		rds.DBs[i] = NewDBInstance(server, i)
	}
	return rds
}

func (rds *RedisDBs) CloseAllDB() {
	for _, db := range rds.DBs {
		db.CloseDB()
	}
}
