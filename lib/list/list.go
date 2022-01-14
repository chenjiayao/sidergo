package list

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

func (l *List) Head() *Node {
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
