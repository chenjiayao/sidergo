package validate

import (
	"fmt"

	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/redis"
)

func ValidateTtl(conn conn.Conn, args [][]byte) error {

	if len(args) != 1 {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command", redis.Ttl)
	}
	return nil
}
