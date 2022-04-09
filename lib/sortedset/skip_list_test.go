package sortedset

import (
	"testing"
)

func TestSkipList_insert(t *testing.T) {
	skipList := MakeSkipList()
	skipList.Print()

	skipList.insert(36, "6")
	skipList.Print()

	skipList.insert(3, "1")
	skipList.Print()

	skipList.insert(12, "3")
	skipList.Print()

	skipList.insert(19, "4")
	skipList.Print()

	skipList.insert(8, "2")
	skipList.Print()

	skipList.insert(23, "5")
	skipList.Print()

}

func TestSkipList_GetRank(t *testing.T) {
	skipList := MakeSkipList()
	skipList.insert(36, "5")
	skipList.insert(3, "0")
	skipList.insert(12, "2")
	skipList.insert(19, "3")
	skipList.insert(8, "1")
	skipList.insert(23, "3")
	skipList.Print()
	got1 := skipList.GetRank("5", 36)
	want1 := 5
	if got1 != int64(want1) {
		t.Errorf("SkipList.GetRank() = %v, want %v", got1, want1)
	}
}
