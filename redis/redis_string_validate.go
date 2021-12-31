package redis

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/chenjiayao/goredistraning/helper"
	"github.com/chenjiayao/goredistraning/rediserr"
)

//set key value [EX seconds] [PX milliseconds] [NX|XX]
func ValidateSet(args [][]byte) error {

	if len(args) < 2 || len(args) > 7 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", set)
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

func ValidateSetNx(args [][]byte) error {
	if len(args) != 2 {
		return rediserr.SYNTAX_ERROR
	}
	return nil
}

func ValidateSetEx(args [][]byte) error {
	if len(args) != 2 {
		return rediserr.SYNTAX_ERROR
	}
	return nil
}

func ValidatePSetEx(args [][]byte) error {
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

func ValidateMSet(args [][]byte) error {
	if len(args)/2 != 0 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", mset)
	}
	return nil
}
func ValidateMGet(args [][]byte) error {
	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", mget)
	}
	return nil
}
func ValidateGetSet(args [][]byte) error {
	if len(args) != 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", getset)
	}
	return nil
}

func ValidateGet(args [][]byte) error {
	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", get)
	}
	return nil
}

func ValidateIncr(args [][]byte) error {
	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", incr)
	}
	return nil
}

func ValidateIncrBy(args [][]byte) error {
	if len(args) != 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", incrby)
	}

	increment := string(args[1])

	_, err := strconv.Atoi(increment)
	if err != nil {
		return rediserr.NOT_INTEGER_ERROR
	}
	return nil
}

func ValidateIncreByFloat(args [][]byte) error {
	if len(args) != 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", incrbyf)
	}

	increment := string(args[1])
	_, err := strconv.ParseFloat(increment, 64)
	if err != nil {
		return rediserr.NOT_INTEGER_ERROR
	}
	return nil
}

func ValidateDecr(args [][]byte) error {
	return ValidateIncr(args)
}

func ValidateDecrBy(args [][]byte) error {
	return ValidateIncrBy(args)
}
