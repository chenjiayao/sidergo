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

func ValidateUnwatch(conn conn.Conn, args [][]byte) error {
	if len(args) > 0 {
		return errors.New("ERR wrong number of arguments for 'unwatch' command")
	}
	return nil
}
