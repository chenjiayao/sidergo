package validate

import (
	"fmt"

	"github.com/chenjiayao/goredistraning/redis"
)

func ValidateSadd(args [][]byte) error {

	if len(args) < 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Set)
	}

	return nil
}

func ValidateSmembers(args [][]byte) error {

	if len(args) > 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Smembers)
	}
	return nil
}

func ValidateScard(args [][]byte) error {
	if len(args) > 2 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Scard)
	}
	return nil
}
