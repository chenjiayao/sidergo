package redis

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chenjiayao/goredistraning/interface/response"
)

var _ response.Response = RedisErrorResponse{}

const (
	CRLF = "\r\n"
)

// 错误：以"-" 开始，如："-ERR Invalid Synatx\r\n"
type RedisErrorResponse struct {
	Err error
}

func (rer RedisErrorResponse) ToContentByte() []byte {
	return []byte{}
}

func (rer RedisErrorResponse) ToErrorByte() []byte {
	// TODO感觉可以优化，后续进行优化
	errString := "-" + rer.Err.Error() + CRLF
	return []byte(errString)
}

func MakeErrorResponse(err string) response.Response {
	return RedisErrorResponse{
		Err: errors.New(err),
	}
}

///////简单字符串：以"+" 开始， 如："+OK\r\n"
type RedisSimpleResponse struct {
	Content string
}

func (rsr RedisSimpleResponse) ToContentByte() []byte {
	content := "+" + rsr.Content + CRLF
	return []byte(content)
}

func (rsr RedisSimpleResponse) ToErrorByte() []byte {
	return []byte{}
}

func MakeSimpleResponse(content [][]byte) response.Response {

	var builder strings.Builder
	for _, b := range content {
		builder.Write(append(b, ' '))
	}

	return RedisSimpleResponse{
		Content: builder.String(),
	}
}

/////整数：以":"开始，如：":1\r\n"
type RedisNumberResponse struct {
	Number int
}

func (rsr RedisNumberResponse) ToContentByte() []byte {
	content := fmt.Sprintf(":%d%s", rsr.Number, CRLF)
	return []byte(content)
}

func (rsr RedisNumberResponse) ToErrorByte() []byte {
	return []byte{}
}

func MakeNumberResponse(number int) response.Response {
	return RedisNumberResponse{
		Number: number,
	}
}
