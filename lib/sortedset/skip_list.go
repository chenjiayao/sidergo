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
	level []*Level
}

type Element struct {
	Member string
	Score  float64
}

func (sl *SkipList) insert(member string, score float64) *Node {
	return nil
}

func (sl *SkipList) remove(member string, score float64) bool {

	//删除某个节点之后，需要更新 next 的指针
	needUpdateNextPointNode := make([]*Node, MAX_LEVEL)
	node := sl.header

	for i := sl.level - 1; i >= 0; i-- {
		for node.level[i].next != nil && (node.level[i].next.Score < score || (node.level[i].next.Score == score && node.level[i].next.Member < member)) {
			node = node.level[i].next
		}
		needUpdateNextPointNode[i] = node
	}

	maybeNeedDelNode := node.level[0].next

	if maybeNeedDelNode != nil && maybeNeedDelNode.Member == member && maybeNeedDelNode.Score == score {
		sl.removeNode(maybeNeedDelNode, needUpdateNextPointNode)
		return true
	}
	return false
}

func (sl *SkipList) removeNode(delNode *Node, nodes []*Node) {
	for i := 0; i < sl.level; i++ {
		if nodes[i].level[i].next == delNode {
			nodes[i].level[i].next = delNode.level[i].next
			nodes[i].level[i].span += delNode.level[i].span - 1
		} else {
			nodes[i].level[i].span--
		}
	}

	if delNode.level[0].next == nil {
		sl.tail = delNode.prev
	} else {
		//处理双向链表
		delNode.level[0].next.prev = delNode.prev
	}
	for sl.level > 1 && sl.header.level[sl.level-1].next == nil {
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
		level: make([]*Level, level),
	}
	for i := range n.level {
		n.level[i] = new(Level)
	}
	return n
}
