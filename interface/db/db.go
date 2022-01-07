package db

import (
	"github.com/chenjiayao/goredistraning/interface/response"
)

type DB interface {
	Exec(cmdName string, args [][]byte) response.Response
}
