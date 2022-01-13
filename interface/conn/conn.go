package conn

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
}
