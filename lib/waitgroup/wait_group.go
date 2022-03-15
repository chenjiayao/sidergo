package waitgroup

import (
	"sync"
	"time"
)

type WaitGroup struct {
	wg sync.WaitGroup
}

func (wg *WaitGroup) WaitWithTimeout(timeout time.Duration) bool {

	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.wg.Wait()
	}()

	select {
	case <-c:
		return false
	case <-time.After(timeout):
		wg.wg.Done()
		return true
	}
}

func (wg *WaitGroup) Add(inc int) {
	wg.wg.Add(inc)
}

func (wg *WaitGroup) Done() {
	wg.wg.Done()
}
