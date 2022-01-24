package conn

import (
	"github.com/chenjiayao/goredistraning/interface/response"
)

type Conn interface {
	Close()
	RemoteAddress() string
	Write(data []byte) error

	GetSelectedDBIndex() int
	SetSelectedDBIndex(index int)

	GetPassword() string
	SetPassword(password string)

	IsInMultiState() bool
	SetMultiState(state int)
	GetMultiState() int

	PushMultiCmd(cmd [][]byte)
	GetMultiCmds() [][][]byte

	Discard()

	DirtyCAS(flag bool)
	GetDirtyCAS() bool

	GetBlockingResponse() response.Response

	SetBlockingResponse(content response.Response)

	SetMaxBlockTime(timeout int64)

	GetMaxBlockTime() int64

	GetBlockingExec() (string, [][]byte)
	SetBlockingExec(cmdName string, args [][]byte)
}
