package redis

import (
	"net"

	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/response"
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

	multiState MultiState //是否处于事务状态

	multiCmdQueues [][][]byte // 事务命令

	redisDirtyCAS bool //标记当前事务是否被破坏 ----> watch 的 key 是否被更改了

	maxBlockTime int64

	blockChan chan response.Response

	blockingCmdName string
	blockingCmdArgs [][]byte
}

func MakeRedisConn(conn net.Conn) *RedisConn {

	rc := &RedisConn{
		conn:           conn,
		selectedDB:     0,
		password:       "",
		multiCmdQueues: make([][][]byte, 0),
		blockChan:      make(chan response.Response),
	}
	return rc
}

func (rc *RedisConn) DirtyCAS(flag bool) {

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

func (rc *RedisConn) GetBlockingResponse() response.Response {
	return <-rc.blockChan
}

func (rc *RedisConn) SetBlockingResponse(resp response.Response) {
	rc.blockChan <- resp
}

func (rc *RedisConn) SetMaxBlockTime(timeout int64) {
	rc.maxBlockTime = timeout
}

func (rc *RedisConn) GetMaxBlockTime() int64 {
	return 0
}

func (rc *RedisConn) GetBlockingExec() (string, [][]byte) {
	return rc.blockingCmdName, rc.blockingCmdArgs
}

func (rc *RedisConn) SetBlockingExec(cmdName string, args [][]byte) {
	rc.blockingCmdName = cmdName
	rc.blockingCmdArgs = args
}
