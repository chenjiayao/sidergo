package resp

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chenjiayao/sidergo/interface/response"
)

var _ response.Response = RedisErrorResponse{}

const (
	CRLF = "\r\n"
)

var (
	NullMultiResponse  = MakeMultiResponse("")
	OKSimpleResponse   = MakeSimpleResponse("OK")
	EmptyArrayResponse = MakeArrayResponse(nil)
)

// 错误：以"-" 开始，如："-ERR Invalid Synatx\r\n"
type RedisErrorResponse struct {
	Err error
}

func (rer RedisErrorResponse) ToContentByte() []byte {
	return []byte{}
}

func (rer RedisErrorResponse) ToStrings() string {
	return rer.Err.Error()
}

func (rer RedisErrorResponse) ToErrorByte() []byte {
	errString := fmt.Sprintf("-%s%s", rer.Err.Error(), CRLF)
	return []byte(errString)
}

func (rer RedisErrorResponse) ISOK() bool {
	return false
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

func (rsr RedisSimpleResponse) ISOK() bool {
	return true
}

func MakeSimpleResponse(content string) response.Response {

	return RedisSimpleResponse{
		Content: content,
	}

}

//////多行数据 $ 开头
/**

$8
codehole

*/
type RedisMultiLineResponse struct {
	Content string
}

func (rmls *RedisMultiLineResponse) ToContentByte() []byte {
	if rmls.Content == "" {
		return []byte("$-1\r\n")
	}

	prefix := fmt.Sprintf("$%d", len(rmls.Content))

	var buffer strings.Builder
	buffer.Write([]byte(prefix))
	buffer.WriteString(CRLF)
	buffer.WriteString(rmls.Content)
	buffer.WriteString(CRLF)
	return []byte(buffer.String())
}

func (rmls *RedisMultiLineResponse) ToErrorByte() []byte {
	return []byte{}
}

func (rmls RedisMultiLineResponse) ISOK() bool {
	return true
}

func MakeMultiResponse(content string) response.Response {

	return &RedisMultiLineResponse{
		Content: content,
	}
}

/////整数：以":"开始，如：":1\r\n"
type RedisNumberResponse struct {
	Number int64
}

func (rsr RedisNumberResponse) ToContentByte() []byte {
	content := fmt.Sprintf(":%d%s", rsr.Number, CRLF)
	return []byte(content)
}

func (rsr RedisNumberResponse) ToErrorByte() []byte {
	return []byte{}
}
func (rsr RedisNumberResponse) ISOK() bool {
	return true
}

func MakeNumberResponse(number int64) response.Response {
	return RedisNumberResponse{
		Number: number,
	}
}

///////
type RedisArrayResponse struct {
	Content []response.Response
}

//*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
func (rar RedisArrayResponse) ToContentByte() []byte {

	if rar.Content == nil {
		return []byte("*0\r\n")
	}

	res := make([]byte, 0)
	res = append(res, []byte(fmt.Sprintf("*%d%s", len(rar.Content), CRLF))...)

	for _, v := range rar.Content {
		res = append(res, v.ToContentByte()...)
	}
	return res
}

func (rar RedisArrayResponse) ToErrorByte() []byte {
	return []byte{}
}
func (rar RedisArrayResponse) ISOK() bool {
	return true
}

func MakeArrayResponse(resps []response.Response) response.Response {
	return RedisArrayResponse{
		Content: resps,
	}
}

/////
type ReidsRawByteResponse struct {
	ByteContent []byte
}

func (rrr ReidsRawByteResponse) ToContentByte() []byte {
	return rrr.ByteContent
}

func (rrr ReidsRawByteResponse) ToErrorByte() []byte {
	return []byte{}
}
func (rrr ReidsRawByteResponse) ISOK() bool {
	return true
}

func MakeReidsRawByteResponse(b []byte) response.Response {
	return ReidsRawByteResponse{
		ByteContent: b,
	}
}
