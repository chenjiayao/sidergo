package list

type List struct {
	head *Node
	tail *Node
	size int64
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

func (node *Node) Prev() *Node {
	return node.prev
}

func (l *List) InsertHead(val interface{}) {
	n := &Node{
		val:  val,
		prev: nil,
		next: l.head,
	}
	if l.head != nil {
		l.head.prev = n
	} else {
		l.tail = n
	}

	l.head = n
	l.size++

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

func (l *List) Len() int64 {
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
	return headNode.Element()
}

func (l *List) InsertBeforePiovt(pivot interface{}, val interface{}) int64 {
	pivotNode := l.getNodeByElement(pivot)
	if pivotNode == nil {
		return -1
	}

	node := &Node{
		val:  val,
		next: pivotNode,
		prev: pivotNode.prev,
	}
	pivotNode.prev = node
	l.size++
	return l.Len()
}

func (l *List) InsertAfterPiovt(pivot interface{}, val interface{}) int64 {
	pivotNode := l.getNodeByElement(pivot)
	if pivotNode == nil {
		return -1
	}

	node := &Node{
		val:  val,
		next: pivotNode.next,
		prev: pivotNode,
	}
	pivotNode.next = node
	l.size++
	return l.Len()
}

func (l *List) getNodeByElement(pivot interface{}) *Node {

	node := l.head
	for {
		if node == nil {
			return nil
		}
		if node.Element() == pivot {
			return node
		}
		node = node.Next()
	}

}

func (l *List) getNodeByIndex(index int64) *Node {
	if l.head == nil {
		return nil
	}

	stop := index
	if index < 0 {
		stop = l.Len() + index
	}

	from := l.head
	for i := 0; ; {
		if int64(i) == stop {
			break
		}
		i++
		from = from.Next()
		if from == nil {
			return nil
		}
	}
	return from
}

// start from 0
// 因为 redis lindex 可以返回 nil，所以 GetElementByIndex 可以返回 nil
func (l *List) GetElementByIndex(index int64) interface{} {
	node := l.getNodeByIndex(index)
	if node == nil {
		return nil
	}
	return node.Element()
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

func (l *List) Trim(start, stop int64) {

	if stop < 0 {
		stop = l.Len() + stop
	}
	if start < 0 {
		start = l.Len() + start
	}

	if stop < start {
		*l = *MakeList()
		return
	}

	startNode := l.getNodeByIndex(start)
	stopNode := l.getNodeByIndex(stop)

	if stop > l.Len() {
		stop = l.Len() - 1
	}
	if start > l.Len() {
		start = l.Len() - 1
	}
	l.size = stop - start + 1
	l.head = startNode
	l.tail = stopNode
}

func (l *List) First() *Node {
	return l.head
}

func (l *List) Last() *Node {
	return l.tail
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

// lrange 可以返回空数组，所以这了只能返回数组
func (l *List) Range(start, stop int64) []interface{} {
	hits := make([]interface{}, 0)

	if start < 0 {
		start = l.Len() + start
	}

	if stop < 0 {
		stop = l.Len() + stop
	}

	if stop < start {
		return hits
	}

	startNode := l.getNodeByIndex(start)
	stopNode := l.getNodeByIndex(stop)

	for {

		if startNode == nil {
			return hits
		}

		hits = append(hits, startNode.Element())
		if startNode == stopNode {
			break
		}
		startNode = startNode.Next()
	}
	return hits
}

func MakeList() *List {

	l := &List{
		head: nil,
		tail: nil,
		size: 0,
	}

	return l
}
