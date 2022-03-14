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
	Set     = "set"
	Setnx   = "setnx"
	Setex   = "setex"
	Psetex  = "psetex"
	Mset    = "mset"
	Mget    = "mget"
	Msetnx  = "msetnx"
	Get     = "get"
	Getset  = "getset"
	Incr    = "incr"
	Incrby  = "incrby"
	Incrbyf = "incrbyfloat"
	Decr    = "decr"
	Decrby  = "decrby"

	//list
	Lpush     = "lpush"
	Lpushx    = "lpushx"
	Rpush     = "rpush"
	Rpushx    = "rpushx"
	Lpop      = "lpop"
	Rpop      = "rpop"
	Ltrim     = "ltrim"
	Rpoplpush = "rpoplpush"
	Lrem      = "lrem"
	Llen      = "llen"
	Lindex    = "lindex"
	Lset      = "lset"
	Linsert   = "linsert"
	Lrange    = "lrange"
	Blpop     = "blpop"
	Brpop     = "brpop"

	//common
	Expire = "expire"
	Del    = "del"
	Rename = "rename"

	//set
	Sadd      = "sadd"
	Smembers  = "smembers"
	Scard     = "scard"
	Spop      = "spop"
	Sismember = "sismember"
	Sdiff     = "sdiff"

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

	Multi   = "multi"
	Discard = "discard"
	Watch   = "watch"
	Exec    = "exec"
	Auth    = "auth"
	Select  = "select"
	Ttl     = "ttl"
	Persist = "Persist"
	Exist   = "Exist"

	Ping = "ping"
)

var (
	CommandTables = make(map[string]Command)

	//写命令，用来判断 aof log
	WriteCommands = map[string]string{
		Set:       "",
		Setnx:     "",
		Setex:     "",
		Persist:   "",
		Psetex:    "",
		Expire:    "",
		Del:       "",
		Lpush:     "",
		Lpushx:    "",
		Rpush:     "",
		Rpushx:    "",
		Lpop:      "",
		Rpop:      "",
		Ltrim:     "",
		Rpoplpush: "",
		Lrem:      "",
		Lindex:    "",
		Lset:      "",
		Linsert:   "",
		Blpop:     "",
		Brpop:     "",

		Mset:    "",
		Msetnx:  "",
		Getset:  "",
		Incr:    "",
		Incrby:  "",
		Incrbyf: "",
		Decr:    "",
		Spop:    "",
		Sadd:    "",
		Decrby:  "",
		Rename:  "",
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
