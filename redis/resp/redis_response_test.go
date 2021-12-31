package resp

import (
	"testing"
)

func TestRedisArrayResponse_ToContentByte(t *testing.T) {

	want := "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"
	b := make([][]byte, 0)
	b = append(b, []byte("SET"))
	b = append(b, []byte("key"))
	b = append(b, []byte("value"))
	res := RedisArrayResponse{
		Content: b,
	}
	t.Run("set key value", func(t *testing.T) {
		got := string(res.ToContentByte())
		if got != want {
			t.Errorf("ToContentByte().toString() = %s, want %s", got, want)
		}
	})

	want = "*2\r\n$3\r\nGET\r\n$7\r\ntestkey\r\n"
	b1 := make([][]byte, 0)
	b1 = append(b1, []byte("GET"))
	b1 = append(b1, []byte("testkey"))
	res = RedisArrayResponse{
		Content: b1,
	}
	t.Run("GET testkey", func(t *testing.T) {
		got := string(res.ToContentByte())
		if got != want {
			t.Errorf("ToContentByte().toString() = %s, want %s", got, want)
		}
	})
}
