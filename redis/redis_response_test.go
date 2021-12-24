package redis

import (
	"bytes"
	"testing"
)

func TestRedisArrayResponse_ToContentByte(t *testing.T) {

	want := []byte("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")
	b := make([][]byte, 0)
	b = append(b, []byte("SET"))
	b = append(b, []byte("key"))
	b = append(b, []byte("value"))
	res := RedisArrayResponse{
		Content: b,
	}
	t.Run("set key value", func(t *testing.T) {
		got := res.ToContentByte()
		if !bytes.Equal(got, want) {
			t.Errorf("ToContentByte() = %v, want %v", got, want)
		}
	})

	want = []byte("*2\r\n$3\r\nGET\r\n$7\r\ntestkey\r\n")
	b1 := make([][]byte, 0)
	b1 = append(b1, []byte("GET"))
	b1 = append(b1, []byte("testkey"))
	res = RedisArrayResponse{
		Content: b1,
	}
	t.Run("GET testkey", func(t *testing.T) {
		got := string(res.ToContentByte())
		if got != string(want) {
			t.Errorf("ToContentByte() = %s, want %s", got, string(want))
		}
	})
}
