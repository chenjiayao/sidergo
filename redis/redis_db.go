package redis

import (
	"fmt"
	"time"

	"github.com/chenjiayao/goredistraning/config"
	"github.com/chenjiayao/goredistraning/interface/db"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/lib/dict"
)

var _ db.DB = &RedisDB{}

const (
	UnlimitTTL = -1
)

type RedisDB struct {
	dataset *dict.ConcurrentDict
	index   int                  // 数据库 db 编号
	ttlMap  *dict.ConcurrentDict //保存 key 和过期时间之间的关系
}

func NewDBInstance(index int) *RedisDB {
	rd := &RedisDB{
		dataset: dict.NewDict(128),
		index:   index,
		ttlMap:  dict.NewDict(128),
	}
	return rd
}

func (rd *RedisDB) Exec(cmdName string, args [][]byte) response.Response {
	if cmdName == "ttl" {
		return MakeNumberResponse(rd.ttl(args[0]))
	}
	return rd.ExecNormal(cmdName, args)
}

// ttl = -2  key 不存在
// ttl = -1 永久有效
func (rd *RedisDB) ttl(key []byte) int64 {
	res, ok := rd.ttlMap.Get(string(key))
	if !ok {
		return -2
	}
	t, _ := res.(time.Time)
	return int64(time.Until(t))
}

//设置key 的 ttl
/*
	保存到 ttlMap 中的是过期的时间
*/
func (rd *RedisDB) setKeyTtl(key []byte, ttl time.Time) {
	rd.ttlMap.Put(string(key), ttl)
}

func (rd *RedisDB) ExecNormal(cmdName string, args [][]byte) response.Response {
	command, ok := CommandTables[cmdName]
	if !ok {
		return MakeErrorResponse(fmt.Sprintf("ERR unknown command '%s'", cmdName))
	}

	//参数校验
	validate := command.ValidateFunc
	err := validate(args)
	if err != nil {
		return MakeErrorResponse(err.Error())
	}

	//执行命令
	execFunc := command.ExecFunc
	resp := execFunc(rd, args)
	return resp
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
