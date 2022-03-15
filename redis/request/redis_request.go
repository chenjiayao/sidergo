package request

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chenjiayao/sidergo/interface/request"
)

var _ request.Request = &RedisRequet{}

type RedisRequet struct {
	CmdName string
	Args    [][]byte //args 本质上是一个字符串数组，不包含 cmdName
	Err     error    //从 socket 读取数据出错
}

func (rr *RedisRequet) ToStrings() string {

	var builder strings.Builder
	builder.WriteString(rr.CmdName + " ")
	for _, v := range rr.Args {
		builder.Write(append(v, ' '))
	}
	return strings.TrimSpace(builder.String())
}

/*

*2
$6
member
$7
member1

*/
func (rr *RedisRequet) ToByte() []byte {

	var builder strings.Builder

	arrLen := len(rr.Args) + 1
	arrHeader := fmt.Sprintf("*%d\r\n", arrLen)
	builder.WriteString(arrHeader)

	cmd := fmt.Sprintf("$%d\r\n%s\r\n", len([]byte(rr.CmdName)), rr.CmdName)
	builder.WriteString(cmd)

	for i := 0; i < len(rr.Args); i++ {
		itemHeader := fmt.Sprintf("$%d\r\n", len(rr.Args[i]))
		item := fmt.Sprintf("%s%s\r\n", itemHeader, string(rr.Args[i]))
		builder.WriteString(item)
	}

	res := builder.String()
	return []byte(res)
}

func (rr *RedisRequet) GetCmdName() string {
	return rr.CmdName
}

//返回参数，不包括命令
func (rr *RedisRequet) GetArgs() [][]byte {
	return rr.Args
}

func (rr *RedisRequet) GetErr() error {
	return rr.Err
}

var (
	PROTOCOL_ERROR_REQUEST = &RedisRequet{Err: errors.New("protocol error")}
)
