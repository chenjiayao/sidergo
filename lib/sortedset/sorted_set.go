package sortedset

type SortedSet struct {
	skipList *SkipList           //排序方式 order by score asc, member asc
	dict     map[string]*Element //不会出现并发读写的问题  member => *Element
}

func MakeSortedSet() *SortedSet {
	return &SortedSet{
		dict:     make(map[string]*Element),
		skipList: makeSkipList(),
	}
}

func (ss *SortedSet) Len() int64 {
	return int64(len(ss.dict))
}

func (ss *SortedSet) Add(memeber string, score float64) bool {
	element, ok := ss.dict[memeber]
	// dict 中，socre 是否相等都可以执行这个逻辑
	ss.dict[memeber] = &Element{
		Member: memeber,
		Score:  score,
	}

	if !ok {
		//skipList.insert
		return true
	}

	//覆盖旧的 score 和 member
	if element.Score != score {
		// skipList 删掉旧的
		// skipList 增加新的
		return true
	}

	return false
}

func (ss *SortedSet) Get(memeber string) (*Element, bool) {
	element, ok := ss.dict[memeber]
	if ok {
		return element, true
	}
	return nil, false
}

func (ss *SortedSet) Count()
