package datatype

import (
	"testing"

	"github.com/chenjiayao/goredistraning/helper"
	"github.com/chenjiayao/goredistraning/lib/set"
	"github.com/chenjiayao/goredistraning/redis"
)

func TestExecSadd(t *testing.T) {
	db := redis.NewDBInstance(0)

	insertValue := [][]byte{
		[]byte("value1"),
		[]byte("value2"),
		[]byte("value3"),
	}

	ExecSadd(nil, db, append([][]byte{[]byte("key")}, insertValue...))

	i, exist := db.Dataset.Get("key")
	if !exist {
		t.Errorf("execAdd should add key to redis")
	}
	setValue := i.(*set.Set)
	vals := setValue.Members()

	ss := helper.BbyteToSString(vals)
	if ss[0] != "value1" {
		t.Errorf("ss[0] = %s, want = %s", ss[0], "value1")
	}

	if ss[1] != "value2" {
		t.Errorf("ss[0] = %s, want = %s", ss[1], "value2")
	}

	if ss[2] != "value3" {
		t.Errorf("ss[0] = %s, want = %s", ss[2], "value3")
	}
}
