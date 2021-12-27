package conn

type Conn interface {
	Close()
	RemoteAddress() string
}
