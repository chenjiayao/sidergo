package validate

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/redis"
)

func ValidateTtl(conn conn.Conn, args [][]byte) error {

	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.TTL)
	}
	return nil
}

func ValidateExpire(conn conn.Conn, args [][]byte) error {
	if len(args) != 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.EXPIRE)
	}
	_, err := strconv.Atoi(string(args[1]))
	if err != nil {
		return errors.New("ERR value is not an integer or out of range")
	}

	return nil
}

func ValidateDel(conn conn.Conn, args [][]byte) error {
	if len(args) < 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.DEL)
	}
	return nil
}

func ValidateRename(conn conn.Conn, args [][]byte) error {
	if len(args) != 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.RENAME)
	}
	return nil
}
