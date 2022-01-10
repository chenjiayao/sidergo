package list

import (
	"testing"
)

//两个接口值相等仅当它们都是nil值或者它们的动态类型相同并且动态值也根据这个动态类型的＝=操作相等。
func TestList_InsertLast(t *testing.T) {
	l := MakeList()
	l.InsertLast(1)
	l.InsertLast(2)
	l.InsertLast(3)
}

func TestList_Exist(t *testing.T) {
	l := MakeList()
	l.InsertLast(1)
	l.InsertLast(2)
	l.InsertLast(3)

	exist := l.Exist(3)
	if !exist {
		t.Errorf("l.Exist(3) = %v, but want true", exist)
	}

	exist = l.Exist(5)
	if exist {
		t.Errorf("l.Exist(5) = %v, but want false", exist)
	}

	exist = l.Exist(1)
	if !exist {
		t.Errorf("l.Exist(1) = %v, but want true", exist)
	}
}

func TestList_Remove(t *testing.T) {
	l := MakeList()
	l.InsertLast(1)
	l.InsertLast(2)
	l.InsertLast(3)

	exist := l.Exist(2)
	if !exist {
		t.Errorf("l.Exist(2) = %v, but want true", exist)
	}

	l.Remove(2)
	exist = l.Exist(2)
	if exist {
		t.Errorf("l.Remove(2) = falied")
	}

	l.Remove(1)
	exist = l.Exist(1)
	if exist {
		t.Errorf("l.Remove(2) = falied")
	}
}
