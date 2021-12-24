package redis

import (
	"net"
	"sync"

	"github.com/chenjiayao/goredistraning/interface/conn"
)

var _ conn.Conn = &RedisConn{}

//每个连接需要保存的信息
type RedisConn struct {
	Conn       net.Conn
	SelectedDB int
	Password   string
	Mu         sync.Mutex
}

func (rc *RedisConn) Close() {
	rc.Conn.Close()
}
