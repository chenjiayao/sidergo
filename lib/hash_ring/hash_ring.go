package hashring

import (
	"hash/crc32"
	"sort"
	"strconv"
)

/*

1. 将 2^32 的整数连接成一个环
2. 将服务器节点的 IP 地址 hash 到环上
3. redis 的每个 key 也 hash 到环上
4. key 的 hash 值在环上顺时针走，碰到的第一个 key hash 就是对应的 node 服务器
5. 服务器节点不一定为真正节点，一个服务器可以虚拟多个节点，分别 hash 在环上的多个值
*/

type HashRing struct {
	replica int //每个真正的node节点，在环山有几个虚拟映射

	nodeMap map[int]string // key 为环上的 int 值，string 为 ip:port 字符串

	hashedKeys []int //node 已经映射在 hash 环上的点，从小到大排序，方便 hit 函数查找--->空间换时间
}

func MakeHashRing(replic int) *HashRing {

	return &HashRing{
		replica:    3,
		nodeMap:    make(map[int]string),
		hashedKeys: make([]int, 0),
	}

}

// key ===> 返回所以在服务器的 IP:port
func (hash *HashRing) Hit(key string) string {
	hashValue := int(crc32.ChecksumIEEE([]byte(key)))

	index := sort.Search(len(hash.hashedKeys), func(i int) bool {
		return hash.hashedKeys[i] >= hashValue
	})

	if index == len(hash.hashedKeys) {
		index = 0
	}
	return hash.nodeMap[hash.hashedKeys[index]]
}

func (hash *HashRing) AddNode(node string) {

	//每个实际的 node 节点在环上会映射出 replica 个虚拟节点
	for i := 0; i < hash.replica; i++ {
		hashValue := int(crc32.ChecksumIEEE([]byte(strconv.Itoa(i) + node)))
		hash.nodeMap[hashValue] = node

		hash.hashedKeys = append(hash.hashedKeys, hashValue)
	}

	sort.Ints(hash.hashedKeys) //排序
}
