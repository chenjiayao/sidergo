package sortedset

import (
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
	skipList := MakeSkipList()
	skipList.insert(36, "5")
	skipList.insert(3, "0")
	skipList.insert(12, "2")
	skipList.insert(19, "3")
	skipList.insert(8, "1")
	skipList.insert(23, "4")
	skipList.Print()
	got1 := skipList.GetRank("5", 36)
	want1 := 5
	if got1 != int64(want1) {
		t.Errorf("SkipList.GetRank() = %v, want %v", got1, want1)
	}
}

func TestSkipList_GetRank(t *testing.T) {
	type fields struct {
		tail   *Node
		header *Node
		level  int
		length int64
	}
	type args struct {
		member string
		score  float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skipList := &SkipList{
				tail:   tt.fields.tail,
				header: tt.fields.header,
				level:  tt.fields.level,
				length: tt.fields.length,
			}
			if got := skipList.GetRank(tt.args.member, tt.args.score); got != tt.want {
				t.Errorf("SkipList.GetRank() = %v, want %v", got, tt.want)
			}
		})
	}
}
