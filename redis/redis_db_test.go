package redis

import (
	"testing"
	"time"

	"github.com/chenjiayao/goredistraning/lib/dict"
)

func TestRedisDB_ttl(t *testing.T) {
	db := &RedisDB{
		dataset: dict.NewDict(6),
		index:   0,
		ttlMap:  dict.NewDict(6),
	}

	args := [][]byte{
		[]byte("key"),
		[]byte("value"),
	}
	ExecSet(db, args)

	time.Sleep(3 * time.Second)
	got := db.ttl([]byte("key"))
	if got != -1 {
		t.Errorf("set key  ttl = -1, but got = %d", got)
	}
}

func TestRedisDB_setKeyTtl(t *testing.T) {
	db := &RedisDB{
		dataset: dict.NewDict(1),
		index:   0,
		ttlMap:  dict.NewDict(1),
	}

	ExecSet(db, [][]byte{
		[]byte("key"),
		[]byte("value"),
	})

	key := []byte("key")
	db.setKeyTtl(key, int64(5*time.Second))

	time.Sleep(3 * time.Second)

	got := db.ttl(key)

	if got != 2 {
		t.Errorf("db.ttl = %d, want= %d", got, 2)
	}

}
