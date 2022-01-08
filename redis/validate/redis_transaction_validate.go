package validate

import (
	"errors"

	"github.com/chenjiayao/goredistraning/interface/conn"
)

func ValidateMultiFun(conn conn.Conn, args [][]byte) error {
	if conn.IsInMultiState() {
		return errors.New("ERR MULTI calls can not be nested")
	}
	return nil
}
