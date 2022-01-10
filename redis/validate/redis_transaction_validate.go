package validate

import (
	"errors"

	"github.com/chenjiayao/goredistraning/interface/conn"
)

func ValidateMulti(conn conn.Conn, args [][]byte) error {
	if conn.IsInMultiState() {
		return errors.New("ERR MULTI calls can not be nested")
	}
	return nil
}

func ValidateDiscard(conn conn.Conn, args [][]byte) error {
	if conn.IsInMultiState() {
		return errors.New("ERR MULTI calls can not be nested")
	}
	return nil
}

func ValidateWatch(conn conn.Conn, args [][]byte) error {
	if conn.IsInMultiState() {
		return errors.New("ERR MULTI calls can not be nested")
	}
	return nil
}
