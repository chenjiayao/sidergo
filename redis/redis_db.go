package redis

import (
	"fmt"

	"github.com/chenjiayao/goredistraning/interface/db"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/lib/dict"
)

var _ db.DB = &RedisDB{}

type RedisDB struct {
	dataset *dict.ConcurrentDict
	index   int // 数据库 db 编号
}

func NewDBInstance(index int) *RedisDB {
	rd := &RedisDB{
		dataset: dict.NewDict(128),
		index:   index,
	}
	return rd
}

func (rd *RedisDB) Exec(cmdName string, args [][]byte) response.Response {
	return rd.ExecNormal(cmdName, args)
}

func (rd *RedisDB) ExecNormal(cmdName string, args [][]byte) response.Response {
	command, ok := CommandTables[cmdName]
	if !ok {
		return MakeErrorResponse(fmt.Sprintf("ERR unknown command '%s'", cmdName))
	}

	//参数校验
	validate := command.ValidateFunc
	ok = validate(args)
	if !ok {
		return MakeErrorResponse(fmt.Sprintf("ERR wrong number of arguments for '%s' command", cmdName))
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
	dbCount := 16
	rds := &RedisDBs{
		DBs:     make([]*RedisDB, dbCount),
		DBCount: dbCount,
	}

	for i := 0; i < dbCount; i++ {
		rds.DBs[i] = NewDBInstance(i)
	}
	return rds
}
