package validate

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/chenjiayao/goredistraning/config"
	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/redis"
)

func ValidateAuthFunc(con conn.Conn, args [][]byte) error {
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

func ValidateSelectFunc(conn conn.Conn, args [][]byte) error {
	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Select)
	}
	dbIndexStr := string(args[0])
	_, err := strconv.Atoi(dbIndexStr)
	if err != nil {
		return errors.New("ERR invalid DB index")
	}
	return nil
}
