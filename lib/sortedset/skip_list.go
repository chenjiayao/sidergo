package sortedset

import (
	"math/rand"
	"time"
)

const (
	MAX_LEVEL = 3
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
	span    int64 // 跳过多少个元素，如果两个元素相邻，那么前一个节点的 span 为 1
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
	updateForwardNodes := make([]*Node, MAX_LEVEL) // 插入新的节点之后，需要更新 forward 指针的节点

	node := skipList.header //node节点最终会定位到「被插入位置之前」

	for i := skipList.level - 1; i >= 0; i-- {
		for node.levels[i].forward != nil && (node.levels[i].forward.Score < score || (node.levels[i].forward.Score == score && node.levels[i].forward.Memeber < memeber)) {
			node = node.levels[i].forward
		}
		updateForwardNodes[i] = node
	}
	levelForNewNode := skipList.RandomLevel()
	newNode := MakeNode(levelForNewNode, score, memeber)

	/**
	newNode 的 levels 会被分成两个部分
	1. levelForNewNode > node.level 那么 node.levels 的每个forward 都指向 newNode，剩余的由更早的 node来指向
	2. levelForNewNode <= node.level 那么 node.levels 从0 到 levelForNewNode 的 level 需要指向newNode
	*/
	if len(newNode.levels) <= len(node.levels) {
		for i := len(newNode.levels) - 1; i >= 0; i-- {
			newNode.levels[i].forward = node.levels[i].forward
			node.levels[i].forward = newNode
		}
	} else {
		for i := len(newNode.levels) - 1; i >= 0; i-- {
			newNode.levels[i].forward = node.levels[i].forward
			node.levels[i].forward = newNode
		}
		//剩余的由更早的 node 来指向，所以需要一个 updateNodes 来保存更早的 node,但是那些更早的 node 只需要更新部分 level
		for i := skipList.level - 1; i <= len(node.levels); i-- {
			newNode.levels[i].forward = updateForwardNodes[i].levels[i].forward
			updateForwardNodes[i].levels[i].forward = newNode
		}
	}

	//插入的新元素是最后一个
	if newNode.levels[0].forward == nil {
		skipList.tail = newNode
	} else {
		node.levels[0].forward.backward = newNode
	}

	//插入的元素是第一个
	if node == skipList.header {
		newNode.backward = nil
	} else {
		newNode.backward = node
	}

	///////更新 span
	/**
	要更新的 span 分成两个部分
	1. skipList.levels ~ newNode.levels 这部分只要自增就行
	2. newNode.levels ~ 1 这部分执行「原来的 span」 - 「newNodes 到下一个节点的 span」+ 1
	*/
	for i := skipList.level - 1; i >= len(newNode.levels); i-- {
		updateForwardNodes[i].levels[i].span++
	}

	for i := len(newNode.levels) - 1; i > 0; i++ {
		updateForwardNodes[i].levels[i].span = updateForwardNodes[i].levels[i].span - newNode.levels[i].span + 1
	}

	skipList.length++
	skipList.reCalculateMaxLevel()

	return newNode
}

//重新计算 skipList 的最大 level
func (skipList *SkipList) reCalculateMaxLevel() {
	for skipList.header.levels[skipList.level-1].forward == nil && skipList.level > 1 {
		skipList.level--
	}
}

func (skipList *SkipList) remove(score float64, member string) *Node {

	updateNodes := make([]*Node, MAX_LEVEL)

	backwardDelNode := skipList.header // node 的下一个节点就是要被删除的节点
	for i := skipList.level - 1; i >= 0; i-- {
		for backwardDelNode.levels[i].forward != nil && (backwardDelNode.levels[i].forward.Score < score || (backwardDelNode.levels[i].forward.Score == score && backwardDelNode.levels[i].forward.Memeber < member)) {
			backwardDelNode = backwardDelNode.levels[i].forward
		}
		updateNodes[i] = backwardDelNode
	}
	removeNode := backwardDelNode.levels[0].forward

	backwardDelNode.levels[0].forward = removeNode.levels[0].forward
	if skipList.tail != removeNode { //删除的不是最后一个元素
		removeNode.backward = backwardDelNode
	}

	if len(backwardDelNode.levels) >= len(removeNode.levels) {
		for i := 0; i < len(backwardDelNode.levels); i++ {
			backwardDelNode.levels[i].forward = removeNode.levels[i].forward
		}
	} else {

		for i := 0; i < len(backwardDelNode.levels); i++ {
			backwardDelNode.levels[i].forward = removeNode.levels[i].forward
		}

		for i := len(backwardDelNode.levels); i < len(removeNode.levels)-1; i++ {
			updateNodes[i].levels[i].forward = removeNode.levels[i].forward
		}
	}

	skipList.length--

	//重新获取最高的 level
	skipList.reCalculateMaxLevel()

	return removeNode
}

func (skipList *SkipList) RandomLevel() int {
	level := 1
	rand.Seed(time.Now().UnixNano())
	for rand.Float32() < P && level < MAX_LEVEL {
		level = level + 1
	}
	return level
}

//如果没有找到，那么返回 -1
func (skipList *SkipList) GetRank(member string, score float64) int64 {
	span := int64(0)
	node := skipList.header

	for i := skipList.level - 1; i >= 0; i-- {
		for node.levels[i].forward != nil && (node.levels[i].forward.Score < score || (node.levels[i].forward.Score == score && node.levels[i].forward.Memeber < member)) {
			span += node.levels[i].span
			node = node.levels[i].forward
		}

		if node.levels[i].forward.Memeber == member {
			span += node.levels[i].span
			return span
		}
	}
	return -1
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
		header: MakeNode(MAX_LEVEL, 0, ""),
		level:  1,
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
