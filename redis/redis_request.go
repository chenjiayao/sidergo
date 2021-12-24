package redis

import (
	"errors"
	"strings"
)

type RedisRequet struct {
	Args [][]byte //args 本质上是一个字符串数组
	Err  error    //从 socket 读取数据出错
}

func (rr RedisRequet) ToStrings() string {

	var builder strings.Builder
	for _, v := range rr.Args {
		builder.Write(append(v, ' '))
	}
	return strings.TrimSpace(builder.String())
}

var (
	PROTOCOL_ERROR_REQUEST = RedisRequet{Err: errors.New("protocol error")}
)
