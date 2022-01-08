package conn

import "github.com/chenjiayao/goredistraning/interface/response"

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

	PushMultiCmd(cmd [][]byte)
	ExecMultiCmds()

	Exec(cmdName string, args [][]byte) response.Response

	Discard()
}
