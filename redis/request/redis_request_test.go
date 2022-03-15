package request

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

func TestRedisRequet_ToByte(t *testing.T) {
	cmd := "*3\r\n$3\r\nget\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"
	want := []byte(cmd)

	req := &RedisRequet{
		CmdName: "get",
		Args: [][]byte{
			[]byte("key"),
			[]byte("value"),
		},
	}
	got := req.ToByte()
	if !SliceEqual(want, got) {
		t.Errorf("test failed")
	}
}

func SliceEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
