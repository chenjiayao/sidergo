package server

import "net"

type Server interface {
	Handle(conn net.Conn)
	Close() error
	Log()
}
