package db

import (
	"github.com/chenjiayao/goredistraning/interface/db"
	"github.com/chenjiayao/goredistraning/lib/dict"
)

var _ db.DB = &RedisDB{}

type RedisDB struct {
	dataset *dict.ConcurrentDict
	index   int // 数据库 db
}

func NewDBInstance(index int) *RedisDB {
	rd := &RedisDB{
		dataset: dict.NewDict(128),
		index:   index,
	}
	return rd
}

func (rd *RedisDB) Exec() {

}

////////////////
type RedisDBs struct {
	db      []*RedisDB
	dbCount int
}

func NewDBs() *RedisDBs {
	dbCount := 16
	rds := &RedisDBs{
		db:      make([]*RedisDB, dbCount),
		dbCount: dbCount,
	}

	for i := 0; i < dbCount; i++ {
		rds.db[i] = NewDBInstance(i)
	}
	return rds
}
