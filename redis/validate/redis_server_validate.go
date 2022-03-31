package validate

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/chenjiayao/sidergo/config"
	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/redis"
)

func ValidateAuth(con conn.Conn, args [][]byte) error {
	if len(args) != 1 {
		return errors.New("ERR wrong number of arguments for 'auth' command")
	}

	if config.Config.RequirePass == "" {
		return errors.New("ERR Client sent AUTH, but no password is set")
	}

	if len(args) != 1 {
		return errors.New("ERR wrong number of arguments for 'auth' command")
	}
	return nil
}

func ValidateSelect(conn conn.Conn, args [][]byte) error {
	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.SELECT)
	}
	dbIndexStr := string(args[0])
	_, err := strconv.Atoi(dbIndexStr)
	if err != nil {
		return errors.New("ERR invalid DB index")
	}
	return nil
}

func ValidatePersist(conn conn.Conn, args [][]byte) error {
	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.PERSIST)
	}
	return nil
}

func ValidateExist(conn conn.Conn, args [][]byte) error {
	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.EXIST)
	}
	return nil
}

func ValidatePing(conn conn.Conn, args [][]byte) error {
	if len(args) > 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.PING)
	}
	return nil
}
