package datatype

import (
	"testing"

	"github.com/chenjiayao/sidergo/lib/list"
	"github.com/chenjiayao/sidergo/redis"
)

func TestExecLrem(t *testing.T) {

}

func TestExecLPush(t *testing.T) {
	db := redis.NewDBInstance(nil, 1)

	args := [][]byte{
		[]byte("list"),
		[]byte("A"),
		[]byte("B"),
	}
	ExecLPush(nil, db, args) //BA
	data, _ := db.Dataset.Get("list")
	listData := data.(*list.List)

	gotl := listData.Len()
	wantl := 2
	if gotl != int64(wantl) {
		t.Errorf("got %d, want %d", gotl, wantl)
	}

	i := listData.PopFromHead()
	if i != "B" {
		t.Errorf("got %s, want %s", i, "B")
	}

	i = listData.PopFromHead()
	if i != "A" {
		t.Errorf("got %s, want %s", i, "A")
	}
}

func TestExecLinsert(t *testing.T) {
	db := redis.NewDBInstance(nil, 1)
	args := [][]byte{
		[]byte("list"),
		[]byte("A"),
		[]byte("C"),
	}
	ExecLPush(nil, db, args) //CA

	//执行成功的结果应该是 ABC
	ExecLinsert(nil, db, [][]byte{
		[]byte("list"),
		[]byte("after"),
		[]byte("A"),
		[]byte("B"), //CAB
	})

	data, _ := db.Dataset.Get("list")
	listData := data.(*list.List)
	gotl := listData.Len()
	wantl := 3
	if gotl != int64(wantl) {
		t.Errorf("got %d, want %d", gotl, wantl)
	}

	i := listData.PopFromHead()
	if i != "C" {
		t.Errorf("got %s, want %s", i, "C")
	}

	i = listData.PopFromHead()
	if i != "A" {
		t.Errorf("got %s, want %s", i, "A")
	}

	i = listData.PopFromHead()
	if i != "B" {
		t.Errorf("got %s, want %s", i, "B")
	}
}

func TestExecLIndex(t *testing.T) {
	db := redis.NewDBInstance(nil, 1)
	args := [][]byte{
		[]byte("list"),
		[]byte("A"),
		[]byte("B"),
		[]byte("C"),
	}
	ExecLPush(nil, db, args) //CBA

	/////
	resp := ExecLIndex(nil, db, [][]byte{
		[]byte("list"),
		[]byte("0"),
	})
	gotStr := string(resp.ToContentByte())
	wantStr := "$1\r\nC\r\n"
	if wantStr != gotStr {
		t.Errorf("got %s, want %s", gotStr, wantStr)
	}

	////
	resp = ExecLIndex(nil, db, [][]byte{
		[]byte("list"),
		[]byte("2"),
	})
	gotStr = string(resp.ToContentByte())
	wantStr = "$1\r\nA\r\n"
	if wantStr != gotStr {
		t.Errorf("got %s, want %s", gotStr, wantStr)
	}

	/////
	resp = ExecLIndex(nil, db, [][]byte{
		[]byte("list"),
		[]byte("1"),
	})
	gotStr = string(resp.ToContentByte())
	wantStr = "$1\r\nB\r\n"
	if wantStr != gotStr {
		t.Errorf("got %s, want %s", gotStr, wantStr)
	}

	/////
	resp = ExecLIndex(nil, db, [][]byte{
		[]byte("list"),
		[]byte("-1"),
	})
	gotStr = string(resp.ToContentByte())
	wantStr = "$1\r\nA\r\n"
	if wantStr != gotStr {
		t.Errorf("got %s, want %s", gotStr, wantStr)
	}

	///
	resp = ExecLIndex(nil, db, [][]byte{
		[]byte("list"),
		[]byte("10"),
	})
	gotStr = string(resp.ToContentByte())
	wantStr = "$-1\r\n"
	if wantStr != gotStr {
		t.Errorf("got %s, want %s", gotStr, wantStr)
	}
}

func TestExecLPop(t *testing.T) {
	db := redis.NewDBInstance(nil, 1)

	list := list.MakeList()
	list.InsertHead("A")
	list.InsertHead("B")
	list.InsertHead("C")
	list.InsertHead("D") //DCBA

	db.Dataset.Put("list", list)
	resp := ExecLPop(nil, db, [][]byte{
		[]byte("list"),
	})
	gotStr := string(resp.ToContentByte())
	wantStr := "$1\r\nD\r\n"
	if wantStr != gotStr {
		t.Errorf("got %s, want %s", gotStr, wantStr)
	}

	resp = ExecLPop(nil, db, [][]byte{
		[]byte("list"),
	})
	gotStr = string(resp.ToContentByte())
	wantStr = "$1\r\nC\r\n"
	if wantStr != gotStr {
		t.Errorf("got %s, want %s", gotStr, wantStr)
	}

}

func TestExecLLen(t *testing.T) {
	db := redis.NewDBInstance(nil, 1)

	args := [][]byte{
		[]byte("list"),
		[]byte("A"),
		[]byte("B"),
	}
	ExecLPush(nil, db, args)

	resp := ExecLLen(nil, db, [][]byte{
		[]byte("list"),
	})
	gotStr := string(resp.ToContentByte())
	wantStr := ":2\r\n"
	if condition := wantStr != gotStr; condition {
		t.Errorf("got %s, want %s", gotStr, wantStr)
	}
}
