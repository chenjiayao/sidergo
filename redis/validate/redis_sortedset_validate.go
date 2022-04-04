package validate

import (
	"errors"
	"strconv"
	"strings"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/lib/border"
)

func ValidateZadd(conn conn.Conn, args [][]byte) error {
	if len(args) < 2 || len(args)/2 == 0 {
		return errors.New("ERR wrong number of arguments for 'zadd' command")
	}

	for i := 1; i < len(args); i += 2 {
		scoreValue := string(args[i])
		_, err := strconv.ParseFloat(scoreValue, 64)
		if err != nil {
			return errors.New("ERR value is not a valid float")
		}
	}
	return nil
}

func ValidateZcard(conn conn.Conn, args [][]byte) error {
	if len(args) != 1 {
		return errors.New("ERR wrong number of arguments for 'zcard' command")
	}
	return nil
}

func ValidateZrank(conn conn.Conn, args [][]byte) error {
	if len(args) != 2 {
		return errors.New("ERR wrong number of arguments for 'zrank' command")
	}
	return nil
}
func ValidateZrevrank(conn conn.Conn, args [][]byte) error {
	if len(args) != 2 {
		return errors.New("ERR wrong number of arguments for 'zrevrank' command")
	}
	return nil
}

func ValidateZscore(conn conn.Conn, args [][]byte) error {
	if len(args) != 2 {
		return errors.New("ERR wrong number of arguments for 'zscore' command")
	}
	return nil
}

func ValidateZrem(conn conn.Conn, args [][]byte) error {
	if len(args) < 2 {
		return errors.New("ERR wrong number of arguments for 'zrem' command")
	}
	return nil
}

func ValidateZremrangebyrank(conn conn.Conn, args [][]byte) error {
	if len(args) != 3 {
		return errors.New("ERR wrong number of arguments for 'zremrangebyrank' command")
	}

	startValue := string(args[1])
	stopValue := string(args[2])

	_, err := strconv.ParseInt(startValue, 10, 64)
	if err != nil {
		return errors.New("ERR value is not an integer or out of range")
	}

	_, err = strconv.ParseInt(stopValue, 10, 64)
	if err != nil {
		return errors.New("ERR value is not an integer or out of range")
	}

	return nil
}

func ValidateZremrangebyscore(conn conn.Conn, args [][]byte) error {
	if len(args) != 3 {
		return errors.New("ERR wrong number of arguments for 'zremrangebyscore' command")
	}

	_, err := border.ParserBorder(string(args[1]))
	if err != nil {
		return err
	}

	_, err = border.ParserBorder(string(args[2]))
	if err != nil {
		return err
	}
	return nil
}

func ValidateZcount(conn conn.Conn, args [][]byte) error {
	if len(args) != 3 {
		return errors.New("ERR wrong number of arguments for 'zcount' command")
	}

	_, err := border.ParserBorder(string(args[1]))
	if err != nil {
		return err
	}

	_, err = border.ParserBorder(string(args[2]))
	if err != nil {
		return err
	}
	return nil
}

func ValidateZincrby(conn conn.Conn, args [][]byte) error {
	if len(args) != 3 {
		return errors.New("ERR wrong number of arguments for 'zincrby' command")
	}
	incrementValue := string(args[1])
	_, err := strconv.ParseFloat(incrementValue, 64)
	if err != nil {
		return errors.New("ERR value is not a valid float")
	}
	return nil
}

func ValidateZrange(conn conn.Conn, args [][]byte) error {
	if len(args) != 3 && len(args) != 4 {
		return errors.New("ERR wrong number of arguments for 'zrange' command")
	}

	startValue := string(args[1])
	stopValue := string(args[2])

	_, err := strconv.ParseInt(startValue, 10, 64)
	if err != nil {
		return errors.New("ERR value is not an integer or out of range")
	}

	_, err = strconv.ParseInt(stopValue, 10, 64)
	if err != nil {
		return errors.New("ERR value is not an integer or out of range")
	}

	if len(args) == 4 && strings.ToLower(string(args[3])) != "withscores" {
		return errors.New("ERR syntax error")
	}

	return nil
}

func ValidateZrevrange(conn conn.Conn, args [][]byte) error {
	if len(args) != 3 && len(args) != 4 {
		return errors.New("ERR wrong number of arguments for 'zrevrange' command")
	}

	startValue := string(args[1])
	stopValue := string(args[2])

	_, err := strconv.ParseInt(startValue, 10, 64)
	if err != nil {
		return errors.New("ERR value is not an integer or out of range")
	}

	_, err = strconv.ParseInt(stopValue, 10, 64)
	if err != nil {
		return errors.New("ERR value is not an integer or out of range")
	}

	if len(args) == 4 && strings.ToLower(string(args[3])) != "withscores" {
		return errors.New("ERR syntax error")
	}

	return nil
}

//ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]
func ValidateZrangebyscore(conn conn.Conn, args [][]byte) error {
	if !(len(args) == 3 || len(args) == 6) {
		return errors.New("ERR wrong number of arguments for 'zrangebyscore' command")
	}
	minValue := string(args[1])
	maxValue := string(args[2])

	_, err := border.ParserBorder(minValue)
	if err != nil {
		return errors.New("ERR value is not an integer or out of range")
	}
	_, err = border.ParserBorder(maxValue)
	if err != nil {
		return errors.New("ERR value is not an integer or out of range")
	}

	if len(args) == 3 {
		return nil
	}

	if strings.ToLower(string(args[3])) != "withscores" {
		return errors.New("ERR syntax error")
	}
	if strings.ToLower(string(args[4])) != "limit" {
		return errors.New("ERR syntax error")
	}

	_, err = strconv.ParseInt(string(args[5]), 10, 64)
	if err != nil {
		return errors.New("ERR value is not an integer or out of range")
	}

	_, err = strconv.ParseInt(string(args[6]), 10, 64)
	if err != nil {
		return errors.New("ERR value is not an integer or out of range")
	}
	return nil
}
