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

func (rr *RedisRequet) ToByte() []byte {
	res := make([]byte, 0)

	for _, v := range rr.Args {
		res = append(res, v...)
		res = append(res, []byte("\r\n")...)
	}
	return res
}

func (rr *RedisRequet) GetCmdName() string {
	return string(rr.Args[0])
}

func (rr *RedisRequet) GetArgs() [][]byte {
	return rr.Args[1:]
}

func (rr *RedisRequet) GetErr() error {
	return rr.Err
}

var (
	PROTOCOL_ERROR_REQUEST = &RedisRequet{Err: errors.New("protocol error")}
)
