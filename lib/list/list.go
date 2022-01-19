package list

import (
	"math"
)

type List struct {
	head *Node
	tail *Node
	size int
}

/**
 *
 *
 *
 */

type Node struct {
	val  interface{}
	prev *Node
	next *Node
}

func (node *Node) Element() interface{} {
	return node.val
}

func (node *Node) Next() *Node {
	return node.next
}

func (l *List) InsertHead(val interface{}) {
	n := &Node{
		val:  val,
		prev: nil,
		next: l.head,
	}
	l.head = n
}

func (l *List) InsertLast(val interface{}) {

	n := &Node{
		val:  val,
		next: nil,
		prev: l.tail,
	}

	if l.head == nil {
		l.head = n
	} else {
		l.tail.next = n
	}

	l.tail = n

	l.size++
}

func (l *List) Len() int {
	return l.size
}

func (l *List) InsertIfNotExist(val interface{}) {
	if !l.Exist(val) {
		l.InsertLast(val)
	}
}

func (l *List) PopFromHead() interface{} {
	if l.head == nil {
		return nil
	}

	headNode := l.head
	l.head = l.head.next
	return headNode
}

func (l *List) GetElementByIndex(index int) interface{} {
	if l.head == nil {
		return nil
	}

	if index == 0 {
		return l.head.val
	}

	var from *Node

	var fromTail bool
	if index >= 0 {
		from = l.head
		fromTail = false
	} else {
		from = l.tail
		fromTail = true
		index = int(math.Abs(float64(index)))
	}

	for i := 1; i <= index; i++ {
		if fromTail {
			from = from.prev
		} else {
			from = from.next
		}

		if from == nil {
			return nil
		}
	}
	return from.val
}

func (l *List) Remove(val interface{}) {
	currentNode := l.head
	for {
		if currentNode == nil {
			return
		}

		if currentNode.Element() == val {
			if currentNode == l.head {
				l.head = currentNode.next
			} else {
				currentNode.prev.next = currentNode.next
			}
		}

		currentNode = currentNode.next
	}
}

func (l *List) First() *Node {
	return l.head
}

func (l *List) Exist(val interface{}) bool {

	exist := false
	cur := l.head

	for {
		if cur == nil {
			break
		}
		if val == cur.val {
			exist = true
			break
		}
		cur = cur.next
	}
	return exist
}

func MakeList() *List {

	l := &List{
		head: nil,
		tail: nil,
		size: 0,
	}

	return l
}
