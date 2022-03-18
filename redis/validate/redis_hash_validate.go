package validate

import (
	"fmt"
	"strconv"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/redis"
)

func ValidateHget(conn conn.Conn, args [][]byte) error {

	if len(args) != 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.HGET)
	}

	return nil
}
func ValidateHmget(conn conn.Conn, args [][]byte) error {

	if len(args) < 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.HGET)
	}

	return nil
}

func ValidateHmset(conn conn.Conn, args [][]byte) error {

	if len(args) < 3 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.HGET)
	}

	return nil
}

func ValidateHset(conn conn.Conn, args [][]byte) error {

	if len(args) != 3 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.HSET)
	}

	return nil
}

func ValidateHsetnx(conn conn.Conn, args [][]byte) error {

	if len(args) != 3 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.HSETNX)
	}

	return nil
}

func ValidateHdel(conn conn.Conn, args [][]byte) error {

	if len(args) < 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.HDEL)
	}
	return nil
}

func ValidateHexists(conn conn.Conn, args [][]byte) error {

	if len(args) != 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.HEXISTS)
	}
	return nil
}

func ValidateHgetall(conn conn.Conn, args [][]byte) error {

	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.HGETALL)
	}

	return nil
}

func ValidateHvals(conn conn.Conn, args [][]byte) error {

	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.HVALS)
	}

	return nil
}

func ValidateHkeys(conn conn.Conn, args [][]byte) error {

	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.HKEYS)
	}

	return nil
}

func ValidateHlen(conn conn.Conn, args [][]byte) error {

	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.HLEN)
	}

	return nil
}

func ValidateHincrby(conn conn.Conn, args [][]byte) error {

	if len(args) != 3 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.HINCRBY)
	}

	increment := string(args[2])
	_, err := strconv.Atoi(increment)
	if err != nil {
		return fmt.Errorf("ERR value is not an integer or out of range")
	}
	return nil
}

func ValidateHincrbyfloat(conn conn.Conn, args [][]byte) error {

	if len(args) != 3 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.HINCRBYFLOAT)
	}

	increment := string(args[2])
	_, err := strconv.ParseFloat(increment, 64)
	if err != nil {
		return fmt.Errorf("ERR value is not an integer or out of range")
	}
	return nil
}
