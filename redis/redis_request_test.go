package redis

import (
	"testing"
)

func TestRedisRequet_ToStrings(t *testing.T) {

	b := make([][]byte, 0)
	b = append(b, []byte("GET"))
	b = append(b, []byte("key"))

	req := &RedisRequet{
		Args: b,
		Err:  nil,
	}
	got := req.ToStrings()
	want := "GET key"

	if got != want {
		t.Errorf("req.ToStrings() = %s, want : %s", got, want)
	}
}
