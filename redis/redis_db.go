package redis

import (
	"strings"

	"github.com/chenjiayao/goredistraning/interface/db"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/lib/dict"
	"github.com/chenjiayao/goredistraning/lib/logger"
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

func (rd *RedisDB) Exec(cmd [][]byte) response.Response {
	cmdName := rd.parseCommand(cmd)
	logger.Info(cmdName)
	return MakeNumberResponse(1)
}

func (rd *RedisDB) parseCommand(cmd [][]byte) string {
	cmdName := string(cmd[0])
	return strings.ToLower(cmdName)
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
