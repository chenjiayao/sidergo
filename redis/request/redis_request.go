package request

import (
	"errors"
	"strings"

	"github.com/chenjiayao/sidergo/interface/request"
)

var _ request.Request = &RedisRequet{}

type RedisRequet struct {
	Args [][]byte //args 本质上是一个字符串数组
	Err  error    //从 socket 读取数据出错
}

func (rr *RedisRequet) ToStrings() string {

	var builder strings.Builder
	for _, v := range rr.Args {
		builder.Write(append(v, ' '))
	}
	return strings.TrimSpace(builder.String())
}

func (rr *RedisRequet) GetArgs() [][]byte {
	return rr.Args
}

var (
	PROTOCOL_ERROR_REQUEST = RedisRequet{Err: errors.New("protocol error")}
)
