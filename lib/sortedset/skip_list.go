package sortedset

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	MAX_LEVEL = 4
)

//跳跃表， sorted set 底层实现
// http://zhangtielei.com/posts/blog-redis-skiplist.html

type Element struct {
	Score  float64
	Member string
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

// skipList  的排序规则为：score, member asc
type SkipList struct {
	tail   *Node
	header *Node
	level  int
	length int64
}

func (skipList *SkipList) insert(score float64, member string) *Node {
	updateForwardNodes := make([]*Node, MAX_LEVEL) // 插入新的节点之后，需要更新 forward 指针的节点

	node := skipList.header //node节点最终会定位到「被插入位置之前」

	for i := MAX_LEVEL - 1; i >= 0; i-- {
		for node.levels[i].forward != nil && (node.levels[i].forward.Score < score || (node.levels[i].forward.Score == score && node.levels[i].forward.Member < member)) {
			node = node.levels[i].forward
		}
		updateForwardNodes[i] = node
	}
	levelForNewNode := skipList.RandomLevel()
	newNode := MakeNode(levelForNewNode, score, member)

	/**
	newNode 的 levels 会被分成两个部分
	1. levelForNewNode > node.level 那么 node.levels 的每个forward 都指向 newNode，剩余的由更早的 node来指向
	2. levelForNewNode <= node.level 那么 node.levels 从0 到 levelForNewNode 的 level 需要指向newNode
	*/
	if levelForNewNode <= len(node.levels) {
		for i := levelForNewNode - 1; i >= 0; i-- {
			newNode.levels[i].forward = node.levels[i].forward
			node.levels[i].forward = newNode
		}
	} else {
		for i := 0; i < len(node.levels); i++ {
			newNode.levels[i].forward = node.levels[i].forward
			node.levels[i].forward = newNode
		}

		//剩余的由更早的 node 来指向，所以需要一个 updateNodes 来保存更早的 node,但是那些更早的 node 只需要更新部分 level
		for i := len(node.levels); i < levelForNewNode; i++ {
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

	for i := len(newNode.levels) - 1; i > 0; i-- {
		updateForwardNodes[i].levels[i].span = updateForwardNodes[i].levels[i].span - newNode.levels[i].span + 1
	}

	skipList.length++
	skipList.reCalculateMaxLevel()

	return newNode
}

//重新计算 skipList 的最大 level
func (skipList *SkipList) reCalculateMaxLevel() {
	skipList.level = MAX_LEVEL - 1
	for skipList.header.levels[skipList.level].forward == nil && skipList.level > 1 {
		skipList.level--
	}
	skipList.level++
}

func (skipList *SkipList) remove(score float64, member string) *Node {

	updateNodes := make([]*Node, MAX_LEVEL)

	backwardDelNode := skipList.header // node 的下一个节点就是要被删除的节点

	for i := skipList.level - 1; i >= 0; i-- {
		for backwardDelNode.levels[i].forward != nil && (backwardDelNode.levels[i].forward.Score < score || (backwardDelNode.levels[i].forward.Score == score && backwardDelNode.levels[i].forward.Member < member)) {
			backwardDelNode = backwardDelNode.levels[i].forward
		}
		updateNodes[i] = backwardDelNode
	}
	removeNode := backwardDelNode.levels[0].forward

	backwardDelNode.levels[0].forward = removeNode.levels[0].forward
	if skipList.tail != removeNode { //删除的不是最后一个元素
		removeNode.backward = backwardDelNode
	}

	//更新 forward 指针
	if len(backwardDelNode.levels) >= len(removeNode.levels) {
		for i := 0; i < len(removeNode.levels); i++ {
			backwardDelNode.levels[i].forward = removeNode.levels[i].forward
		}
	} else {
		//被删除节点的 level > 前一个节点的 level
		for i := 0; i < len(backwardDelNode.levels); i++ {
			backwardDelNode.levels[i].forward = removeNode.levels[i].forward
		}

		for i := len(backwardDelNode.levels); i < len(removeNode.levels)-1; i++ {
			updateNodes[i].levels[i].forward = removeNode.levels[i].forward
		}
	}

	//更新 span 值
	if len(backwardDelNode.levels) >= len(removeNode.levels) {
		//被删除节点的 level <= 前一个节点的 level
		for i := 1; i < len(removeNode.levels); i++ {
			backwardDelNode.levels[i].span = backwardDelNode.levels[i].span + removeNode.levels[i].span - 1
		}

		for i := len(removeNode.levels); i < len(backwardDelNode.levels); i++ {
			backwardDelNode.levels[i].span--
		}

	} else {
		//被删除节点的 level > 前一个节点的 level
		for i := 1; i < len(removeNode.levels); i++ {
			updateNodes[i].levels[i].span = updateNodes[i].levels[i].span + removeNode.levels[i].span - 1
		}

		for i := len(removeNode.levels); i < len(updateNodes); i++ {
			updateNodes[i].levels[i].span = updateNodes[i].levels[i].span - 1
		}
	}

	skipList.length--

	//重新获取最高的 level
	skipList.reCalculateMaxLevel()

	return removeNode
}

func (skipList *SkipList) RandomLevel() int {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(MAX_LEVEL)
	return r + 1
}

//如果没有找到，那么返回 -1
func (skipList *SkipList) GetRank(member string, score float64) int64 {
	span := int64(0)
	currentNode := skipList.header

	for i := skipList.level - 1; i >= 0; i-- {
		for currentNode.levels[i].forward != nil && (currentNode.levels[i].forward.Score < score || (currentNode.levels[i].forward.Score == score && currentNode.levels[i].forward.Member < member)) {
			span += currentNode.levels[i].span
			currentNode = currentNode.levels[i].forward
		}

		if currentNode.levels[i].forward != nil && currentNode.levels[i].forward.Member == member {
			span += currentNode.levels[i].span
			return span
		}
	}
	return -1
}

func (skipList *SkipList) ForEach(start, stop int64, fun func(*Element) bool) {

	node := skipList.header

	index := int64(0)

	for node.levels[0].forward != nil {
		node = node.levels[0].forward
		if index >= start && index <= stop {
			fun(&node.Element)
		}

		index++

		if index > stop {
			break
		}
	}
}

//打印出 skipList 的结构
func (skiplist *SkipList) Print() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	rows := make([]table.Row, 0)

	cols := make([][]string, 0)

	current := skiplist.header
	for i := 0; i < int(skiplist.length)+1; i++ {

		col := make([]string, 0)
		for j := 0; j < MAX_LEVEL; j++ {
			if j < len(current.levels) {
				val := fmt.Sprintf("%0.1f : %d", current.Element.Score, current.levels[j].span)
				col = append(col, val)
			} else {
				val := "nil"
				col = append(col, val)
			}
		}
		cols = append(cols, col)
		current = current.levels[0].forward
	}

	for i := MAX_LEVEL - 1; i >= 0; i-- {
		row := table.Row{}
		for j := 0; j < len(cols); j++ {
			row = append(row, cols[j][i])
		}
		rows = append(rows, row)
	}

	t.AppendRows(rows)
	t.Render()
}

func MakeSkipList() *SkipList {
	return &SkipList{
		tail:   nil,
		header: MakeNode(MAX_LEVEL, 0, ""),
		level:  1,
		length: 0,
	}
}

func MakeNode(level int, score float64, member string) *Node {

	node := &Node{
		Element: Element{
			Score:  score,
			Member: member,
		},
		levels: make([]*Level, level),
	}
	for i := 0; i < len(node.levels); i++ {
		node.levels[i] = &Level{
			forward: nil,
			span:    1,
		}
	}
	return node
}
