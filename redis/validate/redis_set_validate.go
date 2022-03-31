package validate

import (
	"fmt"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/redis"
)

func ValidateSadd(conn conn.Conn, args [][]byte) error {

	if len(args) < 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.SET)
	}

	return nil
}

func ValidateSmembers(conn conn.Conn, args [][]byte) error {

	if len(args) > 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.SMEMBERS)
	}
	return nil
}

func ValidateScard(conn conn.Conn, args [][]byte) error {
	if len(args) > 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.SCARD)
	}
	return nil
}

func ValidateSpop(conn conn.Conn, args [][]byte) error {
	if len(args) > 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.SPOP)
	}
	return nil
}

func ValidateSismember(conn conn.Conn, args [][]byte) error {
	if len(args) != 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.SISMEMBER)
	}
	return nil
}

func ValidateSmove(conn conn.Conn, args [][]byte) error {
	if len(args) != 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.SMOVE)
	}
	return nil
}

func ValidateSdiff(conn conn.Conn, args [][]byte) error {
	if len(args) < 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.SDIFF)
	}
	return nil
}
