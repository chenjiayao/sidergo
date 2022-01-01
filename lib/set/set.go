package set

type Set struct {
	//TODO 是否直接使用  []byte 当作 key 会不会更高效，这个需要进行压测试试
	vals map[string]struct{}
}

func (set *Set) Add(v string) int {
	_, exist := set.vals[v]
	if exist {
		return 0
	}

	//unsafe.Sizeof(struct{}{}) == 0
	set.vals[v] = struct{}{}
	return 1
}

func (set *Set) Len() int {
	return len(set.vals)
}

func (set *Set) Del(v string) {
	delete(set.vals, v)
}

func (set *Set) Members() [][]byte {
	keys := make([][]byte, len(set.vals))
	for k, v := range set.vals {
		_ = v //避免 simplify range expressionsimplifyrange
		keys = append(keys, []byte(k))
	}
	return keys
}

func MakeSet(size int64) *Set {
	s := &Set{
		vals: make(map[string]struct{}, size),
	}
	return s
}
