package redis

import (
	"net"

	"github.com/chenjiayao/goredistraning/interface/conn"
)

var _ conn.Conn = &RedisConn{}

//每个连接需要保存的信息
type RedisConn struct {
	conn       net.Conn
	selectedDB int
	password   string

	inMultiState   bool       //是否处于事务状态
	multiCmdQueues [][][]byte // 事务命令
}

func MakeRedisConn(conn net.Conn) *RedisConn {

	rc := &RedisConn{
		conn:           conn,
		selectedDB:     0,
		password:       "",
		multiCmdQueues: make([][][]byte, 128),
	}
	return rc
}

func (rc *RedisConn) IsInMultiState() bool {
	return rc.inMultiState
}

func (rc *RedisConn) SetMultiState(state int) {
	if state == 1 {
		rc.inMultiState = true
	} else {
		rc.inMultiState = false
	}
}

func (rc *RedisConn) GetPassword() string {
	return rc.password
}

func (rc *RedisConn) SetPassword(password string) {
	rc.password = password
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

func (rc *RedisConn) SetSelectedDBIndex(index int) {
	rc.selectedDB = index
}
