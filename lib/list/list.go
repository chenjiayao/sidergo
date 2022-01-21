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

func (l *List) getNodeByIndex(index int) *Node {
	if l.head == nil {
		return nil
	}

	stop := index
	if index < 0 {
		stop = l.Len() + index
	}

	from := l.head
	for i := 0; ; {
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
func (l *List) GetElementByIndex(index int) interface{} {
	node := l.getNodeByIndex(index)
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

func (l *List) Range(start int, stop int) []interface{} {
	startNode := l.getNodeByIndex(start)
	stopNode := l.getNodeByIndex(stop)

	hits := make([]interface{}, 0)

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
