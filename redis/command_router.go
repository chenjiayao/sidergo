package redis

import (
	"strings"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/response"
)

type RedisExecCommandFunc func(conn conn.Conn, db *RedisDB, args [][]byte) response.Response
type RedisExecValidateFunc func(conn conn.Conn, args [][]byte) error

const (
	//
	HDEL         = "hdel"
	HEXISTS      = "hexists"
	HGET         = "hget"
	HGETALL      = "hgetall"
	HINCRBY      = "hincrby"
	HINCRBYFLOAT = "hincrbyfloat"
	HKEYS        = "hkeys"
	HLEN         = "hlen"
	HMGET        = "hmget"
	HMSET        = "hmset"
	HSET         = "hset"
	HSETNX       = "hsetnx"
	HVALS        = "hvals"

	//string
	SET     = "set"
	SETNX   = "setnx"
	SETEX   = "setex"
	PSETEX  = "psetex"
	MSET    = "mset"
	MGET    = "mget"
	MSETNX  = "msetnx"
	GET     = "get"
	GETSET  = "getset"
	INCR    = "incr"
	INCRBY  = "incrby"
	INCRBYF = "incrbyfloat"
	DECR    = "decr"
	DECRBY  = "decrby"

	//list
	LPUSH     = "lpush"
	LPUSHX    = "lpushx"
	RPUSH     = "rpush"
	RPUSHX    = "rpushx"
	LPOP      = "lpop"
	RPOP      = "rpop"
	LTRIM     = "ltrim"
	RPOPLPUSH = "rpoplpush"
	LREM      = "lrem"
	LLEN      = "llen"
	LINDEX    = "lindex"
	LSET      = "lset"
	LINSERT   = "linsert"
	LRANGE    = "lrange"
	BLPOP     = "blpop"
	BRPOP     = "brpop"

	//common
	EXPIRE = "expire"
	DEL    = "del"
	RENAME = "rename"

	//set
	SADD      = "sadd"
	SMEMBERS  = "smembers"
	SCARD     = "scard"
	SPOP      = "spop"
	SISMEMBER = "sismember"
	SDIFF     = "sdiff"
	SMOVE     = "smove"

	//sorted set
	ZADD             = "zadd"
	ZCARD            = "zcard"
	ZCOUNT           = "zcount"
	ZINCRBY          = "zincrby"
	ZRANGE           = "zrange"
	ZRANGEBYSCORE    = "zrangebyscore"
	ZRANK            = "zrank"
	ZREM             = "zrem"
	ZREMRANGEBYRANK  = "zremrangebyrank"
	ZREMRANGEBYSCORE = "zremrangebyscore"
	ZREVRANGE        = "zrevrange"
	ZREVRANGEBYSCORE = "zrevrangebyscore"
	ZREVRANK         = "zrevrank"
	ZSCORE           = "zscore"
	ZUNIONSTORE      = "zunionstore"
	ZINTERSTORE      = "zinterstore"
	ZSCAN            = "zscan"

	MULTI   = "multi"
	DISCARD = "discard"
	WATCH   = "watch"
	EXEC    = "exec"
	UNWATCH = "unwatch"
	AUTH    = "auth"
	SELECT  = "select"
	TTL     = "ttl"
	PERSIST = "Persist"
	EXIST   = "Exist"

	PING = "ping"
)

var (
	CommandTables = make(map[string]Command)

	//写命令，用来判断 aof log
	WriteCommands = map[string]string{
		SET:       "",
		SETNX:     "",
		SETEX:     "",
		PERSIST:   "",
		PSETEX:    "",
		EXPIRE:    "",
		DEL:       "",
		LPUSH:     "",
		LPUSHX:    "",
		RPUSH:     "",
		RPUSHX:    "",
		LPOP:      "",
		RPOP:      "",
		LTRIM:     "",
		RPOPLPUSH: "",
		LREM:      "",
		LINDEX:    "",
		LSET:      "",
		LINSERT:   "",
		BLPOP:     "",
		BRPOP:     "",

		MSET:    "",
		MSETNX:  "",
		GETSET:  "",
		INCR:    "",
		INCRBY:  "",
		INCRBYF: "",
		DECR:    "",
		SPOP:    "",
		SADD:    "",
		DECRBY:  "",
		RENAME:  "",
	}
)

type Command struct {
	CmdName      string
	CommandFunc  RedisExecCommandFunc
	ValidateFunc RedisExecValidateFunc
}

func RegisterRedisCommand(cmdName string, commandFunc RedisExecCommandFunc, validateFunc RedisExecValidateFunc) {

	cmdName = strings.ToLower(cmdName)
	CommandTables[cmdName] = Command{
		CmdName:      cmdName,
		CommandFunc:  commandFunc,
		ValidateFunc: validateFunc,
	}
}
