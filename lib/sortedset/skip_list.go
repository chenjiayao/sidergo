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

type Element struct {
	Score   float64
	Memeber string
}

type Level struct {
	forward *Node // 同层的下一个节点
	span    int64 // 跳过多少个元素
}
type Node struct {
	Element
	levels   []*Level // len(levels) 是随机出来的
	backward *Node    //  最底层的前一个节点
}

// skipList  的排序规则为：score, memeber asc
type SkipList struct {
	tail   *Node
	header *Node
	level  int
	length int64
}

func (skipList *SkipList) insert(score float64, memeber string) *Node {
	updateNode := make([]*Node, MAX_LEVEL)

	updateSpan := make([]*Node, MAX_LEVEL)

	node := skipList.header //node节点最终会定位到「被插入位置之前」

	for i := skipList.level - 1; i >= 0; i-- {

		for node.levels[i] != nil && (node.Score < score || (node.Score == score && node.Memeber < memeber)) {
			node = node.levels[i].forward
		}

		updateNode[i] = node
	}
	levelForNewNode := skipList.RandomLevel()
	newNode := MakeNode(levelForNewNode, score, memeber)

	/**
	newNode 的 levels 可能会被分成两个部分
	1. levelForNewNode > node.level 那么 node.levels 的每个forward 都指向 newNode，剩余的由更早的 node来指向
	2. levelForNewNode <= node.level 那么 node.levels 从0 到 levelForNewNode 的 level 需要指向newNode
	*/
	if len(newNode.levels) <= len(node.levels) {
		for i := len(node.levels) - 1; i >= 0; i-- {
			newNode.levels[i].forward = node.levels[i].forward
			node.levels[i].forward = newNode
		}
	} else {
		for i := len(node.levels) - 1; i >= 0; i-- {
			newNode.levels[i].forward = node.levels[i].forward
			node.levels[i].forward = newNode
		}
		//剩余的由更早的 node 来指向，所以需要一个 updateNode 来保存更早的 node,但是那些更早的 node 只需要更新部分 level
		for i := skipList.level - 1; i <= len(node.levels); i-- {
			newNode.levels[i].forward = updateNode[i].levels[i].forward
			updateNode[i].levels[i].forward = newNode
		}
	}

	//插入的新元素是最后一个
	if node == skipList.tail {
		skipList.tail = newNode
	}

	newNode.backward = node

	skipList.length++
	if skipList.level < levelForNewNode {
		skipList.level = levelForNewNode
	}

	return newNode
}

func (skipList *SkipList) RandomLevel() int {
	level := 1
	rand.Seed(time.Now().UnixNano())
	for rand.Float32() < P && level < MAX_LEVEL {
		level = level + 1
	}
	return level
}

func (skipList *SkipList) Find(member string) (float64, bool) {
	node := skipList.header
	for i := skipList.level - 1; i >= 0; i-- {
		for node.levels[i] != nil && node.Memeber < member {
			node = node.levels[i].forward
		}

		if node.Memeber == member {
			return node.Score, true
		}
	}
	return 0, false
}

func MakeSkipList() *SkipList {
	return &SkipList{
		tail:   nil,
		header: MakeNode(0, 0, ""),
		level:  0,
		length: 0,
	}
}

func MakeNode(level int, score float64, memeber string) *Node {

	node := &Node{
		Element: Element{
			Score:   score,
			Memeber: memeber,
		},
		levels: make([]*Level, level),
	}
	for i := 0; i < len(node.levels); i++ {
		node.levels[i] = &Level{
			forward: nil,
			span:    0,
		}
	}
	return node
}
