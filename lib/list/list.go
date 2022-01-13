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

	tail := l.tail

	node := &Node{
		next: nil,
		prev: tail,
		val:  val,
	}
	tail.next = node
	l.tail = node

	l.size++
}

func (l *List) InsertIfNotExist(val interface{}) {
	if !l.Exist(val) {
		l.InsertLast(val)
	}
}

func (l *List) Remove(val interface{}) {
	cur := l.head

	for {
		if cur == nil {
			break
		}

		if val == cur.val {
			cur.prev.next = cur.next
		}

		cur = cur.next
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
	n := &Node{
		val: 0,
	}

	n.prev = nil
	n.next = nil

	l := &List{
		head: n,
		tail: n,
		size: 0,
	}

	return l
}
