package server

import (
	"net"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/request"
	"github.com/chenjiayao/sidergo/interface/response"
)

type Server interface {
	Handle(conn net.Conn)
	Close() error
	Log()
	Exec(client conn.Conn, request request.Request) response.Response
}
