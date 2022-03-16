package request

type Request interface {
	GetArgs() [][]byte
	GetErr() error
	ToByte() []byte
	GetCmdName() string
}
