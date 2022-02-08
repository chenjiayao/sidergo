package sortedset

type SortedSet struct {
	skipList *SkipList
	dict     map[string]*Element
}

func MakeSortedSet() *SortedSet {
	return &SortedSet{
		dict:     make(map[string]*Element),
		skipList: makeSkipList(),
	}
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

	if element.Score != score {
		// skipList 删掉旧的
		// skipList 增加新的
		return true
	}

	return false
}
