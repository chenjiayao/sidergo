package list

type List struct {
	head *Node
	tail *Node
	size int64
}

func (l *List) HeadNode() *Node {
	return l.head
}

func (l *List) TailNode() *Node {
	return l.tail
}

func (l *List) InsertTail(val interface{}) {

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

func (l *List) PopFromHead() interface{} {
	if l.head == nil {
		return nil
	}

	headNode := l.head
	l.head = l.head.next
	l.size--
	return headNode.Element()
}

func (l *List) PopFromTail() interface{} {
	if l.tail == nil {
		return nil
	}

	node := l.tail
	l.tail = l.tail.prev
	l.size--
	return node.Element()
}

func (l *List) Len() int64 {
	return l.size
}

func (l *List) InsertIfNotExist(val interface{}) {
	if !l.Exist(val) {
		l.InsertTail(val)
	}
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

	if pivotNode == l.head {
		l.head = node
	}

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

	if pivotNode == l.tail {
		l.tail = node
	}

	pivotNode.next = node
	l.size++
	return l.Len()
}

func (l *List) getNodeByElement(val interface{}) *Node {
	node := l.head
	for {
		if node == nil {
			return nil
		}
		if node.Element() == val {
			return node
		}
		node = node.Next()
	}

}

func (l *List) GetNodeByIndex(index int64) *Node {
	if l.head == nil {
		return nil
	}

	stop := index
	if index < 0 {
		stop = l.Len() + index
	}

	from := l.head
	for i := int64(0); ; {
		if i == stop {
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
	node := l.GetNodeByIndex(index)
	if node == nil {
		return nil
	}
	return node.Element()
}

func (l *List) RemoveNode(val interface{}) {

	removeNode := l.getNodeByElement(val)
	if removeNode == nil {
		return
	}

	l.size--
	if l.head == removeNode {
		l.head = removeNode.Next()
		removeNode.Next().prev = nil
		return
	}

	if l.tail == removeNode {
		l.tail = removeNode.Prev()
	}

	removeNode.Prev().next = removeNode.Next()
}

//保留[start, stop]
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

	startNode := l.GetNodeByIndex(start)
	stopNode := l.GetNodeByIndex(stop)

	if stop > l.Len() {
		stop = l.Len() - 1
	}
	if start > l.Len() {
		start = l.Len() - 1
	}

	l.size = stop - start + 1
	l.head = startNode
	l.tail = stopNode

	startNode.prev = nil
	stopNode.next = nil
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

	startNode := l.GetNodeByIndex(start)
	stopNode := l.GetNodeByIndex(stop)

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

type Node struct {
	val  interface{}
	prev *Node
	next *Node
}

func (node *Node) Element() interface{} {
	return node.val
}

func (node *Node) SetElement(val interface{}) {
	node.val = val
}

func (node *Node) Next() *Node {
	return node.next
}

func (node *Node) Prev() *Node {
	return node.prev
}
