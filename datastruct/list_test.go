package datastruct

import (
	"testing"
)

func TestList_InsertLast(t *testing.T) {
	l := MakeList()
	l.InsertLast(1)
	l.InsertLast(2)
	l.InsertLast(3)
}
