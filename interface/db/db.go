package db

import (
	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/response"
)

type DB interface {
	Exec(conn conn.Conn, cmdName string, args [][]byte) response.Response
}
