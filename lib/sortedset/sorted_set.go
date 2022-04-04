package sortedset

import "github.com/chenjiayao/sidergo/lib/border"

type SortedSet struct {
	skipList *SkipList           //排序方式 order by score asc, member asc
	dict     map[string]*Element //不会出现并发读写的问题  member => *Element
}

func MakeSortedSet() *SortedSet {
	return &SortedSet{
		dict:     make(map[string]*Element),
		skipList: MakeSkipList(),
	}
}

func (ss *SortedSet) Len() int64 {
	return int64(len(ss.dict))
}

func (ss *SortedSet) Add(member string, score float64) bool {
	element, ok := ss.dict[member]
	// dict 中，socre 是否相等都可以执行这个逻辑
	ss.dict[member] = &Element{
		Member: member,
		Score:  score,
	}

	if !ok {
		//skipList.insert
		ss.skipList.insert(score, member)
		return true
	}

	//覆盖旧的 score 和 member
	if element.Score != score {
		// skipList 删掉旧的
		// skipList 增加新的
		ss.skipList.remove(score, member)
		ss.skipList.insert(score, member)
		return true
	}

	return false
}

func (ss *SortedSet) Get(member string) (*Element, bool) {
	element, ok := ss.dict[member]
	if ok {
		return element, true
	}
	return nil, false
}

func (ss *SortedSet) Count(minBorder, maxBorder *border.Border) int64 {
	i := int64(0)
	return i
}

func (ss *SortedSet) GetRank(member string, score float64) int64 {
	return ss.skipList.GetRank(member, score)
}

func (ss *SortedSet) Remove(member string) bool {
	element, exist := ss.dict[member]
	if exist {
		ss.skipList.remove(element.Score, element.Member)
		delete(ss.dict, member)
		return true
	}

	return false
}

func (ss *SortedSet) Range(start, stop int64) []*Element {
	elements := make([]*Element, 0)
	ss.skipList.ForEach(start, stop, func(e *Element) bool {
		elements = append(elements, e)
		return true
	})
	return elements
}
