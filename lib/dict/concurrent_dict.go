package dict

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const prime32 = uint32(16777619)

//并发安全的 map
// 实现方式：map 分为多段，每段都有一个 lock，减少争锁的可能性
type ConcurrentDict struct {
	fragmentCount int         //分段数
	fragments     []*Fragment // 分段
	count         int32       // map 元素总个数
}

type Fragment struct {
	data map[string]interface{}
	lock sync.RWMutex
}

func (d *ConcurrentDict) RandomKey() interface{} {

	if d.fragmentCount <= 0 {
		return nil
	}

	allKeys := []string{}
	for i := 0; i < d.fragmentCount; i++ {
		fragment := d.getFragment(i)
		for key := range fragment.data {
			allKeys = append(allKeys, key)
		}
	}

	rand.Seed(int64(time.Now().UnixNano()))
	randomIndex := rand.Intn(len(allKeys))
	return allKeys[randomIndex]
}

func (d *ConcurrentDict) Get(key string) (interface{}, bool) {
	if d == nil {
		panic("dict is null")
	}

	hashKey := fnv32(key)
	index := d.spread(hashKey)
	fragment := d.getFragment(int(index))

	fragment.lock.Lock()
	defer fragment.lock.Unlock()

	val, exists := fragment.data[key]
	return val, exists
}

func (d *ConcurrentDict) spread(hashCode uint32) uint32 {
	if d == nil {
		panic("dict is nil")
	}
	tableSize := uint32(len(d.fragments))
	return (tableSize - 1) & uint32(hashCode)
}

//获取分段
func (d *ConcurrentDict) getFragment(index int) *Fragment {
	return d.fragments[index]
}

func (d *ConcurrentDict) Put(key string, val interface{}) bool {
	if d == nil {
		panic("dict is null")
	}

	hashKey := fnv32(key)
	index := d.spread(hashKey)
	fragment := d.getFragment(int(index))

	fragment.lock.Lock()
	defer fragment.lock.Unlock()

	if _, ok := fragment.data[key]; !ok {
		d.increaseCount()
	}
	fragment.data[key] = val
	return true
}

func (d *ConcurrentDict) PutIfExist(key string, val interface{}) bool {
	hashKey := fnv32(key)
	index := d.spread(hashKey)
	fragment := d.getFragment(int(index))

	fragment.lock.Lock()
	defer fragment.lock.Unlock()

	if _, ok := fragment.data[key]; ok {
		fragment.data[key] = val
	}
	return true
}

func (d *ConcurrentDict) PutIfNotExist(key string, val interface{}) bool {
	hashKey := fnv32(key)
	index := d.spread(hashKey)
	fragment := d.getFragment(int(index))

	fragment.lock.Lock()
	defer fragment.lock.Unlock()

	if _, ok := fragment.data[key]; !ok {
		fragment.data[key] = val
		return true
	}
	return false
}

func (d *ConcurrentDict) Range(f func(key string, val interface{})) {
	for i := 0; i < d.fragmentCount; i++ {
		fragment := d.getFragment(i)
		fragment.lock.RLock()
		for key, val := range fragment.data {
			f(key, val)
		}
		fragment.lock.RUnlock()
	}
}

func (d *ConcurrentDict) Del(key string) bool {
	if d == nil {
		panic("dict is null")
	}

	hashKey := fnv32(key)
	index := d.spread(hashKey)
	fragment := d.getFragment(int(index))

	fragment.lock.Lock()
	defer fragment.lock.Unlock()

	delete(fragment.data, key)
	d.decreaseCount()
	return true
}

func (d *ConcurrentDict) Len() int32 {
	return atomic.LoadInt32(&d.count)
}

func (d *ConcurrentDict) increaseCount() {
	atomic.AddInt32(&d.count, 1)
}

func (d *ConcurrentDict) decreaseCount() {
	atomic.AddInt32(&d.count, -1)
}

func (d *ConcurrentDict) Clear() {
	*d = *NewDict(d.fragmentCount)
}

func NewDict(fragmentCount int) *ConcurrentDict {

	d := &ConcurrentDict{
		count:         0,
		fragmentCount: fragmentCount,
		fragments:     make([]*Fragment, fragmentCount),
	}
	for i := 0; i < fragmentCount; i++ {
		d.fragments[i] = &Fragment{
			data: make(map[string]interface{}),
		}
	}
	return d
}

// fnv32 is hash
func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
