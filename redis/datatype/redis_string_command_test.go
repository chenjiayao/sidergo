package datatype

import (
	"bytes"
	"testing"

	"github.com/chenjiayao/sidergo/lib/dict"
	"github.com/chenjiayao/sidergo/redis"
	"github.com/chenjiayao/sidergo/redis/redisresponse"
)

func TestExecSet(t *testing.T) {
	db := redis.NewDBInstance(nil, 0)

	args := [][]byte{
		[]byte("key"),
		[]byte("value"),
	}
	got := ExecSet(nil, db, args)
	want := redisresponse.OKSimpleResponse
	if got != want {
		t.Errorf(" ExecSet(nil,db, args) = %v, want = %v", got, want)
	}

	v, ok := db.Dataset.Get("key")
	if !ok {
		t.Errorf("execSet failed")
	}
	res := v.(string)
	if res != "value" {
		t.Errorf("set store value, but got = %s", res)
	}

	ttl := ttl(db, [][]byte{[]byte("key")})
	if ttl != -1 {
		t.Errorf("set key  ttl = -1, but got = %d", ttl)
	}
}

func TestExecGet(t *testing.T) {

	db := redis.NewDBInstance(nil, 0)
	key := "key"
	value := "value"
	db.Dataset.Put(key, value)

	res := ExecGet(nil, db, [][]byte{
		[]byte(key),
	})

	want := string(redisresponse.MakeSimpleResponse("value").ToContentByte())
	if !bytes.Equal(res.ToContentByte(), []byte(want)) {
		t.Errorf("ExecGet = %s, want %s", string(res.ToContentByte()), want)
	}
}

func TestExecIncrBy(t *testing.T) {
	db := redis.NewDBInstance(nil, 0)

	args := [][]byte{
		[]byte("key"),
		[]byte("1"),
	}
	ExecSet(nil, db, args)

	ExecIncr(nil, db, [][]byte{[]byte("key")})

	v, _ := db.Dataset.Get("key")
	got, _ := v.(string)
	if got != "2" {
		t.Errorf("execIncr should incr key to 2, but key = %s now", got)
	}
}

func TestExecGetset(t *testing.T) {
	db := &redis.RedisDB{
		Dataset: dict.NewDict(6),
		Index:   0,
		TtlMap:  dict.NewDict(6),
	}
	key := "key"
	value := "value"
	db.Dataset.Put(key, value)

	newValue := "newvalue"
	res := ExecGetset(nil, db, [][]byte{
		[]byte("key"),
		[]byte(newValue),
	})
	want := redisresponse.MakeSimpleResponse(value)
	if string(string(want.ToContentByte())) != string(res.ToContentByte()) {
		t.Errorf("execgetSet = %s, want = %s", string(res.ToContentByte()), "+value")
	}
	s, _ := getAsString(nil, db, []byte(key))
	if newValue != s {
		t.Errorf("execgetset store %s , but get %s", "newvalue", s)
	}
}

func TestExecSetEX(t *testing.T) {
	db := redis.NewDBInstance(nil, 0)

	args := [][]byte{
		[]byte("key"),
		[]byte("value"),
	}
	db.Dataset.Put("key", "value")
	resp := ExecSetNX(nil, db, args)
	wantByte := []byte(":0\r\n")
	if !bytes.Equal(resp.ToContentByte(), wantByte) {
		t.Errorf("ExecSetNX = %s, want = %s", string(resp.ToContentByte()), string(wantByte))
	}
}
