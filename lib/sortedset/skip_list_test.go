package sortedset

import (
	"testing"
)

func TestSkipList_insert(t *testing.T) {
	skipList := MakeSkipList()
	skipList.insert(10, "member1")
}
