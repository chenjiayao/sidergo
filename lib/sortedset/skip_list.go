package sortedset

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
	Element
	prev  *Node //前一个节点地址，这个只有最底层的链表才有，最底层的链表是一个双向链表
	Level []*Level
}

type Element struct {
	Member string
	Score  float64
}

func (sl *SkipList) Foreach(f func(element *Element) bool) {
	node := sl.header
	for node != nil {
		f(&node.Element)
		node = node.Level[0].next
	}
}

// insert node
// update span
// update length
// (maybe) update level
//  (maybe) update tail
func (sl *SkipList) insert(member string, score float64) *Node {

	update := make([]*Node, MAX_LEVEL) // link new node with node in `update`
	rank := make([]int64, MAX_LEVEL)

	// find position to insert
	node := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		if i == sl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1] // store rank that is crossed to reach the insert position
		}
		if node.Level[i] != nil {
			// traverse the skip list
			for node.Level[i].next != nil &&
				(node.Level[i].next.Score < score ||
					(node.Level[i].next.Score == score && node.Level[i].next.Member < member)) { // same score, different key
				rank[i] += node.Level[i].span
				node = node.Level[i].next
			}
		}
		update[i] = node
	}

	level := sl.RandomLevel()
	// extend sl level
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.header
			update[i].Level[i].span = int64(sl.length)
		}
		sl.level = level
	}

	// make node and link into sl
	node = makeNode(level, member, score)
	for i := 0; i < level; i++ {
		node.Level[i].next = update[i].Level[i].next
		update[i].Level[i].next = node

		// update span covered by update[i] as node is inserted here
		node.Level[i].span = update[i].Level[i].span - (rank[0] - rank[i])
		update[i].Level[i].span = (rank[0] - rank[i]) + 1
	}

	// increment span for untouched levels
	for i := level; i < sl.level; i++ {
		update[i].Level[i].span++
	}

	// set prev node
	if update[0] == sl.header {
		node.prev = nil
	} else {
		node.prev = update[0]
	}
	if node.Level[0].next != nil {
		node.Level[0].next.prev = node
	} else {
		sl.tail = node
	}
	sl.length++
	return node
}

func (sl *SkipList) remove(member string, score float64) bool {

	//删除某个节点之后，需要更新 next 的指针
	needUpdateNextPointNode := make([]*Node, sl.level)

	node := sl.header

	for i := sl.level - 1; i >= 0; i-- {
		for node.Level[i].next != nil && (node.Level[i].next.Score < score || (node.Level[i].next.Score == score && node.Level[i].next.Member < member)) {
			node = node.Level[i].next
		}
		needUpdateNextPointNode[i] = node
	}

	mayNeedDelNode := node.Level[0].next

	if mayNeedDelNode != nil && mayNeedDelNode.Member == member && mayNeedDelNode.Score == score {
		delNode := mayNeedDelNode
		sl.removeNode(delNode, needUpdateNextPointNode)
		return true
	}
	return false
}

func (sl *SkipList) removeNode(delNode *Node, nodes []*Node) {
	for i := 0; i < sl.level; i++ {
		if nodes[i].Level[i].next == delNode {
			nodes[i].Level[i].next = delNode.Level[i].next
			nodes[i].Level[i].span += delNode.Level[i].span - 1
		} else {
			nodes[i].Level[i].span--
		}
	}

	if delNode.Level[0].next == nil {
		sl.tail = delNode.prev
	} else {
		//处理双向链表
		delNode.Level[0].next.prev = delNode.prev
	}
	for sl.level > 1 && sl.header.Level[sl.level-1].next == nil {
		sl.level--
	}

	sl.length--
}

//根据 rank 获取 node
func (sl *SkipList) GetByRank(rank int64) *Node {

	skipListLevel := sl.level

	node := sl.header

	j := int64(0)

	for i := int64(skipListLevel - 1); i >= 0; i-- {
		for node.Level[i].next != nil && i+node.Level[i].span <= rank {
			node = node.Level[i].next
			j += node.Level[i].span
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

func makeSkipList() *SkipList {
	header := makeNode(16, "", 0)
	return &SkipList{
		level:  0,
		length: 0,
		header: header,
	}
}

func makeNode(level int, member string, score float64) *Node {

	n := &Node{
		Element: Element{
			Member: member,
			Score:  score,
		},
		Level: make([]*Level, level),
	}
	for i := range n.Level {
		n.Level[i] = new(Level)
	}
	return n
}
