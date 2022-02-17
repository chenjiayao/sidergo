package db

import (
	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/response"
)

type DB interface {
	Exec(conn conn.Conn, cmdName string, args [][]byte) response.Response
}
