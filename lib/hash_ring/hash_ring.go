package hashring

type HashRing struct {
	replica int //每个真正的node节点，在环山有几个虚拟映射

	nodeMap map[int]string // key 为环上的 int 值，string 为 ip:port 字符串

}

func MakeHashRing(replic int) *HashRing {

	return nil
}

// key ===> 返回所以在服务器的 IP:port
func (hash *HashRing) Hit(key string) string {
	return ""
}

func (hash *HashRing) AddNodes(nodes ...string) {

}
