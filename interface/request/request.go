package request

type Request interface {
	ToStrings() string
	GetArgs() [][]byte
}
