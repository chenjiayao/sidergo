package validate

import (
	"fmt"
	"strconv"

	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/redis"
)

func ValidateLPush(conn conn.Conn, args [][]byte) error {

	if len(args) < 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Lpush)
	}
	return nil
}

func ValidateLPop(conn conn.Conn, args [][]byte) error {

	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Lpop)
	}
	return nil
}

func ValidateLIndex(conn conn.Conn, args [][]byte) error {
	if len(args) != 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Lindex)
	}
	indexs := args[1]
	_, err := strconv.Atoi(string(indexs))
	if err != nil {
		return fmt.Errorf("(error) ERR value is not an integer or out of range")
	}
	return nil
}

func ValidateLLen(conn conn.Conn, args [][]byte) error {
	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Llen)
	}
	return nil
}

func ValidateLPushx(conn conn.Conn, args [][]byte) error {
	if len(args) != 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Lpush)
	}
	return nil
}
