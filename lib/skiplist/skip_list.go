package skiplist

//跳跃表， sorted set 底层实现

type SkipList struct {
	level  int
	header *Node
	tail   *Node
	length int //最底层链表的长度
}

type Node struct {
	Element,
	NextNodes []*Node //下一个节点，这里 NextNodes[0] 表示最底层的节点，也就是「原始链表」
}

type Element struct {
	Member string
	Score  float64
}

func MakeSkipList() *SkipList {
	header := &Node{
		Element: nil,
		NextNodes: []*Node{
			nil,
		},
	}
	return &SkipList{
		level:  0,
		header: header,
		tail:   nil,
		length: 0,
	}
}
