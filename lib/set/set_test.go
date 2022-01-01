package set

import (
	"testing"
)

func TestSet_Add(t *testing.T) {
	s := MakeSet(128)
	r := s.Add("key")
	if r != 1 {
		t.Errorf("s.add() = %d, but want %d", r, 1)
	}
}

func TestSet_Len(t *testing.T) {
	s := MakeSet(128)
	s.Add("key")

	l := s.Len()
	if l != 1 {
		t.Errorf("s.len = %d, but want = %d", l, 1)
	}
}

func TestSet_Members(t *testing.T) {
	s := MakeSet(128)
	s.Add("key")

	sm := s.Members()
	if string(sm[0]) != "key" {
		t.Errorf("sm[0] = %s, but want = %s", string(sm[0]), "key")
	}
}
