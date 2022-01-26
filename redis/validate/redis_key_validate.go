package validate

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/redis"
)

func ValidateTtl(conn conn.Conn, args [][]byte) error {

	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Ttl)
	}
	return nil
}

func ValidateExpire(conn conn.Conn, args [][]byte) error {
	if len(args) != 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Expire)
	}
	_, err := strconv.Atoi(string(args[1]))
	if err != nil {
		return errors.New("(error) ERR value is not an integer or out of range")
	}

	return nil
}

func ValidateDel(conn conn.Conn, args [][]byte) error {
	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Del)
	}
	return nil
}
