package validate

import "github.com/chenjiayao/sidergo/interface/conn"

func ValidatePrepare(conn conn.Conn, args [][]byte) error {
	return nil
}

func ValidateCommit(conn conn.Conn, args [][]byte) error {
	return nil
}

func ValidateUndo(conn conn.Conn, args [][]byte) error {
	return nil
}
