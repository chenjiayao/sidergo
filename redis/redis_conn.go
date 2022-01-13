package redis

import (
	"net"

	"github.com/chenjiayao/goredistraning/interface/conn"
)

var _ conn.Conn = &RedisConn{}

type MultiState int

const (
	NotInMultiState MultiState = iota
	InMultiState
	InMultiStateButHaveError
)

//每个连接需要保存的信息
type RedisConn struct {
	conn       net.Conn
	selectedDB int
	password   string

	multiState     MultiState
	multiCmdQueues [][][]byte // 事务命令

	redisDirtyCAS bool //标记当前事务是否被破坏 ----> watch 的 key 是否被更改了
}

func MakeRedisConn(conn net.Conn) *RedisConn {

	rc := &RedisConn{
		conn:           conn,
		selectedDB:     0,
		password:       "",
		multiCmdQueues: make([][][]byte, 0),
	}
	return rc
}

func (rc *RedisConn) DirtyCAS(flag bool) {
	if !rc.IsInMultiState() {
		return
	}
	rc.redisDirtyCAS = flag
}

func (rc *RedisConn) GetDirtyCAS() bool {
	return rc.redisDirtyCAS
}

func (rc *RedisConn) Discard() {
	rc.SetMultiState(int(NotInMultiState))
	rc.multiCmdQueues = rc.multiCmdQueues[:0]
}
func (rc *RedisConn) IsInMultiState() bool {
	return rc.multiState == InMultiState || rc.multiState == InMultiStateButHaveError
}

func (rc *RedisConn) SetMultiState(state int) {
	rc.multiState = MultiState(state)
}

func (rc *RedisConn) PushMultiCmd(cmd [][]byte) {
	rc.multiCmdQueues = append(rc.multiCmdQueues, cmd)
}

func (rc *RedisConn) GetMultiCmds() [][][]byte {
	return rc.multiCmdQueues
}

func (rc *RedisConn) GetMultiState() int {
	return int(rc.multiState)
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
