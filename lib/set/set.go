package set

type Set struct {
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

func (set *Set) Exist(key string) bool {
	_, exist := set.vals[key]
	return exist
}

func (set *Set) Del(v string) {
	delete(set.vals, v)
}

func (set *Set) Members() [][]byte {
	sl := len(set.vals)
	keys := make([][]byte, sl)
	i := 0
	for key, value := range set.vals {
		_ = value

		if key == "" {
			continue
		}

		keys[i] = []byte(key)
		i++
	}
	return keys
}

func MakeSet(size int64) *Set {
	s := &Set{
		vals: make(map[string]struct{}, size),
	}
	return s
}
