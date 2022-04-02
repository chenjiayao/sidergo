package sortedset

import (
	"log"
	"testing"
)

func TestSkipList_insert(t *testing.T) {
	skipList := MakeSkipList()
	skipList.insert(36, "6")
	skipList.insert(3, "1")
	skipList.insert(12, "3")
	skipList.insert(19, "4")
	skipList.insert(8, "2")
	skipList.insert(23, "5")

	skipList.Print()
}

func TestSkipList_RandomLevel(t *testing.T) {
	skiplist := MakeSkipList()
	for i := 0; i < 100; i++ {
		r := skiplist.RandomLevel()
		log.Println(r)
	}
}
