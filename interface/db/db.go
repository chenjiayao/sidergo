package db

import "github.com/chenjiayao/goredistraning/interface/response"

type DB interface {
	Exec(cmds [][]byte) response.Response
}
