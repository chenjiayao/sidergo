package sortedset

import "github.com/chenjiayao/goredistraning/lib/border"

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
		ss.skipList.insert(memeber, score)
		return true
	}

	//覆盖旧的 score 和 member
	if element.Score != score {
		// skipList 删掉旧的
		// skipList 增加新的
		ss.skipList.remove(memeber, score)
		ss.skipList.insert(memeber, score)
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

func (ss *SortedSet) Count(minBorder, maxBorder *border.Border) int64 {
	i := int64(0)
	ss.skipList.Foreach(func(element *Element) bool {
		if minBorder.Greater(element.Score) && maxBorder.Less(element.Score) {
			i++
			return true
		}
		return false
	})
	return i
}
