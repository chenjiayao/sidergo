package redis

import (
	"net"
	"sync"

	"github.com/chenjiayao/goredistraning/interface/conn"
)

var _ conn.Conn = &RedisConn{}

//每个连接需要保存的信息
type RedisConn struct {
	conn       net.Conn
	selectedDB int
	Password   string
	Mu         sync.Mutex
}

func MakeRedisConn(conn net.Conn) *RedisConn {
	rc := &RedisConn{
		conn:       conn,
		selectedDB: 0,
		Password:   "",
	}
	return rc
}

func (rc *RedisConn) Close() {
	rc.conn.Close()
}

func (rc *RedisConn) Write(data []byte) error {
	_, err := rc.conn.Write(data)
	return err
}

func (rc *RedisConn) RemoteAddress() string {
	return rc.conn.RemoteAddr().String()
}

func (rc *RedisConn) GetSelectedDBIndex() int {
	return rc.selectedDB
}
