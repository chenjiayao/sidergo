package datastruct

type List struct {
	head *Node
	tail *Node
	size int
}

type Node struct {
	val  interface{}
	prev *Node
	next *Node
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
