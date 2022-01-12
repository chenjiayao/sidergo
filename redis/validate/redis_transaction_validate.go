package validate

import (
	"errors"
	"fmt"

	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/redis"
)

func ValidateMulti(conn conn.Conn, args [][]byte) error {
	if conn.IsInMultiState() {
		return errors.New("ERR MULTI calls can not be nested")
	}
	return nil
}

func ValidateExec(conn conn.Conn, args [][]byte) error {
	if !conn.IsInMultiState() {
		return errors.New("ERR EXEC without MULTI")
	}
	if len(args) > 0 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Exec)
	}
	return nil
}

func ValidateDiscard(conn conn.Conn, args [][]byte) error {
	if !conn.IsInMultiState() {
		return errors.New("ERR DISCARD without MULTI")
	}
	return nil
}

func ValidateWatch(conn conn.Conn, args [][]byte) error {
	if len(args) < 1 {
		return errors.New("ERR wrong number of arguments for 'watch' command")
	}
	return nil
}

func ValidateUnwatch(conn conn.Conn, args [][]byte) error {
	if len(args) > 0 {
		return errors.New("ERR wrong number of arguments for 'unwatch' command")
	}
	return nil
}
