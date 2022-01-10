package validate

import (
	"fmt"

	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/redis"
)

func ValidateSadd(conn conn.Conn, args [][]byte) error {

	if len(args) < 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Set)
	}

	return nil
}

func ValidateSmembers(conn conn.Conn, args [][]byte) error {

	if len(args) > 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Smembers)
	}
	return nil
}

func ValidateScard(conn conn.Conn, args [][]byte) error {
	if len(args) > 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Scard)
	}
	return nil
}

func ValidateSpop(conn conn.Conn, args [][]byte) error {
	if len(args) > 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Spop)
	}
	return nil
}

func ValidateSismember(conn conn.Conn, args [][]byte) error {
	if len(args) != 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Sismember)
	}
	return nil
}

func ValidateSdiff(conn conn.Conn, args [][]byte) error {
	if len(args) < 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Sdiff)
	}
	return nil
}
