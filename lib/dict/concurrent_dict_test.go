package dict

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
)

func TestConcurrentDict_Get(t *testing.T) {

	d := NewDict(6)
	for i := 0; i < 100; i++ {
		d.Put(fmt.Sprintf("test_%d", i), i)
	}
	if d.count != 100 {
		t.Errorf("d.count = %d, want %d", d.count, 100)
	}

	for i := 0; i < 100; i++ {
		val, has := d.Get(fmt.Sprintf("test_%d", i))
		if !has {
			t.Errorf("d.Get should have ,but not")
			continue
		}
		if v, ok := val.(int); !ok || v != i {
			t.Errorf("d.Get = %d  ,want %d", v, i)
		}
	}
}

func TestConcurrentDict_Del(t *testing.T) {
	d := NewDict(6)
	for i := 0; i < 100; i++ {
		d.Put(fmt.Sprintf("test_%d", i), i)
	}

	for i := 0; i < 20; i++ {
		val, _ := rand.Int(rand.Reader, big.NewInt(100))
		key := fmt.Sprintf("test_%d", val)

		d.Del(key)

		_, ok := d.Get(key)
		if ok {
			t.Errorf("d.Del should removed, but not")
		}
	}

	got := d.Len()
	if got != 80 {
		t.Errorf("d len want 80. but got = %d", got)
	}
}

func TestConcurrentDict_Clear(t *testing.T) {
	d := NewDict(6)
	for i := 0; i < 100; i++ {
		d.Put(fmt.Sprintf("test_%d", i), i)
	}
	d.Clear()
	got := d.Len()
	if got != 0 {
		t.Errorf("d len want 0. but got = %d", got)
	}
}
