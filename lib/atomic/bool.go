package atomic

import "sync/atomic"

type Boolean int32

func (b *Boolean) Get() bool {
	return atomic.LoadInt32((*int32)(b)) != 0
}

func (b *Boolean) Set(v bool) {
	if v {
		atomic.StoreInt32((*int32)(b), 1)
	} else {
		atomic.StoreInt32((*int32)(b), 0)
	}
}
