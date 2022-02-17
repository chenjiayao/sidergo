package validate

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/chenjiayao/sidergo/helper"
	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/redis"
	"github.com/chenjiayao/sidergo/redis/rediserr"
)

//set key value [EX seconds] [PX milliseconds] [NX|XX]
func ValidateSet(conn conn.Conn, args [][]byte) error {

	if len(args) < 2 || len(args) > 7 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Set)
	}

	ss := helper.BbyteToSString(args)

	shouldHaveArgsCount := 2

	EXFlagIndex := helper.ContainWithoutCaseSensitive(ss, "EX")
	PXFlagIndex := helper.ContainWithoutCaseSensitive(ss, "PX")
	NXFlagIndex := helper.ContainWithoutCaseSensitive(ss, "NX")
	XXFlagIndex := helper.ContainWithoutCaseSensitive(ss, "XX")

	if NXFlagIndex != -1 && XXFlagIndex != -1 {
		return rediserr.SYNTAX_ERROR
	}
	if PXFlagIndex != -1 && EXFlagIndex != -1 {
		return rediserr.SYNTAX_ERROR
	}

	if EXFlagIndex != -1 {
		ex := ss[EXFlagIndex+1]
		_, err := strconv.Atoi(ex) //ex 下一个参数得是 integer
		if err != nil {
			return errors.New("ERR value is not an integer or out of range")
		}
		shouldHaveArgsCount += 2
	}
	if PXFlagIndex != -1 {
		px := ss[EXFlagIndex+1]
		_, err := strconv.Atoi(px) //px 下一个参数得是 integer
		if err != nil {
			return errors.New("ERR value is not an integer or out of range")
		}
		shouldHaveArgsCount += 2
	}

	if NXFlagIndex != -1 {
		shouldHaveArgsCount += 1
	}
	if XXFlagIndex != -1 {
		shouldHaveArgsCount += 1
	}

	if shouldHaveArgsCount != len(ss) {
		return rediserr.SYNTAX_ERROR
	}
	return nil
}

func ValidateSetNx(conn conn.Conn, args [][]byte) error {
	if len(args) != 2 {
		return rediserr.SYNTAX_ERROR
	}
	return nil
}

func ValidateSetEx(conn conn.Conn, args [][]byte) error {
	if len(args) != 2 {
		return rediserr.SYNTAX_ERROR
	}
	return nil
}

func ValidatePSetEx(conn conn.Conn, args [][]byte) error {
	if len(args) != 3 {
		return rediserr.SYNTAX_ERROR
	}
	mttl := string(args[1])
	_, err := strconv.Atoi(mttl)
	if err != nil {
		return rediserr.NOT_INTEGER_ERROR
	}
	return nil
}

func ValidateMSet(conn conn.Conn, args [][]byte) error {
	if len(args)%2 != 0 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Mset)
	}
	return nil
}

func ValidateMSetNX(conn conn.Conn, args [][]byte) error {
	return ValidateMSet(conn, args)
}

func ValidateMGet(conn conn.Conn, args [][]byte) error {
	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Mget)
	}
	return nil
}
func ValidateGetSet(conn conn.Conn, args [][]byte) error {
	if len(args) != 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Getset)
	}
	return nil
}

func ValidateGet(conn conn.Conn, args [][]byte) error {
	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Get)
	}
	return nil
}

func ValidateIncr(conn conn.Conn, args [][]byte) error {
	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Incr)
	}
	return nil
}

func ValidateIncrBy(conn conn.Conn, args [][]byte) error {
	if len(args) != 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Incrby)
	}

	increment := string(args[1])

	_, err := strconv.Atoi(increment)
	if err != nil {
		return rediserr.NOT_INTEGER_ERROR
	}
	return nil
}

func ValidateIncreByFloat(conn conn.Conn, args [][]byte) error {
	if len(args) != 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Incrbyf)
	}

	increment := string(args[1])
	_, err := strconv.ParseFloat(increment, 64)
	if err != nil {
		return rediserr.NOT_INTEGER_ERROR
	}
	return nil
}

func ValidateDecr(conn conn.Conn, args [][]byte) error {
	return ValidateIncr(conn, args)
}

func ValidateDecrBy(conn conn.Conn, args [][]byte) error {
	return ValidateIncrBy(conn, args)
}
