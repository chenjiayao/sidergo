package validate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/redis"
)

func ValidateLPush(conn conn.Conn, args [][]byte) error {

	if len(args) < 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Lpush)
	}
	return nil
}

func ValidateRPush(conn conn.Conn, args [][]byte) error {

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

func ValidateRPop(conn conn.Conn, args [][]byte) error {
	return ValidateLPop(conn, args)
}

func ValidateLIndex(conn conn.Conn, args [][]byte) error {
	if len(args) != 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Lindex)
	}
	indexs := args[1]
	_, err := strconv.Atoi(string(indexs))
	if err != nil {
		return fmt.Errorf("ERR value is not an integer or out of range")
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

func ValidateRPushx(conn conn.Conn, args [][]byte) error {
	if len(args) != 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Rpushx)
	}
	return nil
}

func ValidateLTrim(conn conn.Conn, args [][]byte) error {
	if len(args) != 3 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Ltrim)
	}
	_, err := strconv.Atoi(string(args[1]))
	if err != nil {
		return fmt.Errorf("ERR value is not an integer or out of range")
	}

	_, err = strconv.Atoi(string(args[2]))
	if err != nil {
		return fmt.Errorf("ERR value is not an integer or out of range")
	}
	return nil
}

func ValidateLrange(conn conn.Conn, args [][]byte) error {
	return ValidateLTrim(conn, args)
}

func ValidateLInsert(conn conn.Conn, args [][]byte) error {
	if len(args) != 4 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Ltrim)
	}
	pos := strings.ToUpper(string(args[1]))
	if pos != "BEFORE" && pos != "AFTER" {
		return fmt.Errorf("ERR syntax error")
	}
	return nil
}

func ValidateLset(conn conn.Conn, args [][]byte) error {
	if len(args) != 3 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Lset)
	}

	_, err := strconv.Atoi(string(args[1]))
	if err != nil {
		return errors.New(" WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	return nil
}

func ValidateBlpop(conn conn.Conn, args [][]byte) error {
	if len(args) < 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Blpop)
	}
	_, err := strconv.Atoi(string(args[len(args)-1]))
	if err != nil {
		return fmt.Errorf("ERR value is not an integer or out of range")
	}
	return nil
}

func ValidateBrpop(conn conn.Conn, args [][]byte) error {
	if len(args) < 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Brpop)
	}
	_, err := strconv.Atoi(string(args[len(args)-1]))
	if err != nil {
		return fmt.Errorf("ERR value is not an integer or out of range")
	}
	return nil
}

func ValidateLrem(conn conn.Conn, args [][]byte) error {
	if len(args) < 3 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Lrem)
	}
	_, err := strconv.Atoi(string(args[1]))
	if err != nil {
		return fmt.Errorf("ERR value is not an integer or out of range")
	}
	return nil
}
