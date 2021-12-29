package redis

import (
	"testing"

	"github.com/chenjiayao/goredistraning/lib/dict"
)

func TestExecSet(t *testing.T) {
	db := &RedisDB{
		dataset: dict.NewDict(6),
		index:   0,
		ttlMap:  dict.NewDict(6),
	}

}
