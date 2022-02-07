package skiplist

import (
	"math/rand"
	"time"
)

const (
	MAX_LEVEL = 32
	P         = 0.25
)

//跳跃表， sorted set 底层实现
// http://zhangtielei.com/posts/blog-redis-skiplist.html
type SkipList struct {
	level  int
	header *Node
	tail   *Node
	length int //最底层链表的长度
}

type Level struct {
	next *Node // 指向同层中的下一个节点
	span int64 // 到 next 跳过的节点数,这个数据用来计算 rank 排名
}

type Node struct {
	Element,
	backward *Node
	level []*Level
}

type Element struct {
	Member string
	Score  float64
}

func MakeSkipList() *SkipList {
	header := &Node{
		Element:  nil,
		backward: nil,
		level:    nil,
	}

	tail := &Node{
		Element:  nil,
		backward: nil,
		level:    nil,
	}

	return &SkipList{
		level:  0,
		length: 0,
		header: header,
		tail:   tail,
	}
}

func (sl *SkipList) Insert(member string, score float64) *Node {
	return nil
}

func (sl *SkipList) GetScoreByMember(member string) *Node {

}

func (sl *SkipList) GetByRank(rank int64) *Node {

	skipListLevel := sl.level

	node := sl.header

	j := int64(0)

	for i := int64(skipListLevel - 1); i >= 0; i-- {
		for node.level[i].next != nil && i+node.level[i].span <= rank {
			node = node.level[i].next
			j += node.level[i].span
		}

		if j == rank {
			return node
		}
	}
	return nil
}

func (sl *SkipList) Del(start int64, stop int64) []*Element {
	return nil
}

func (sl *SkipList) RandomLevel() int {
	level := 1
	rand.Seed(time.Now().UnixNano())
	for rand.Float32() < P && level < MAX_LEVEL {
		level = level + 1
	}
	return level
}
