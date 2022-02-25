package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
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

func main() {
	listener, _ := net.Listen("tcp", ":3101")
	conn, _ := listener.Accept()
	ch := ReadCommand(conn)

	for {
		select {
		case req := <-ch:
			log.Println(req.ToStrings())

			resp := make([]byte, 0)
			header := fmt.Sprintf("*%d\r\n", len(req.Args))
			resp = append(resp, []byte(header)...)

			for _, v := range req.Args { //  set key value
				argLen := fmt.Sprintf("$%d\r\n", len(v))
				resp = append(resp, []byte(argLen)...)
				resp = append(resp, v...)
				resp = append(resp, []byte("\r\n")...)
			}
			conn.Write(resp)
		default:
			continue
		}
	}
}

/*
 为了这个实例代码后续可直接服用到 sidergo 项目中
 这里额外开一个 goroutine 解析，把解析的结果放到 chan 中传递
*/
func ReadCommand(reader net.Conn) chan RedisRequet {
	ch := make(chan RedisRequet)
	go ParseFromSocket(reader, ch)
	return ch
}

func ParseFromSocket(reader io.Reader, ch chan RedisRequet) {
	buf := bufio.NewReader(reader)

	for {
		cmds := [][]byte{}
		header, err := buf.ReadBytes('\n')
		if err != nil {

			ch <- RedisRequet{
				Err: err,
			}
			//如果是客户端关闭了，那么就不要读了，直接退出当前协程
			if io.EOF == err {
				break
			}
			continue
		}

		//首个字符不是 *， 协议错误
		if header[0] != '*' {
			ch <- PROTOCOL_ERROR_REQUEST
			continue
		}
		argsCount, err := parseCmdArgsCount(header)
		if err != nil {
			ch <- PROTOCOL_ERROR_REQUEST
			continue
		}

		//依次读取 数组参数
		// *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
		readArgsFail := false

		for i := 0; i < argsCount; i++ {
			argsWithDelimiter, err := buf.ReadBytes('\n')
			if err != nil {
				ch <- RedisRequet{
					Err: err,
				}
				//如果是客户端关闭了，那么就不要读了，直接退出当前协程
				if io.EOF == err {
					break
				}
				readArgsFail = true
				break
			}

			// $3\r\n
			if argsWithDelimiter[0] != '$' ||
				argsWithDelimiter[len(argsWithDelimiter)-1] != '\n' ||
				argsWithDelimiter[len(argsWithDelimiter)-2] != '\r' {

				ch <- PROTOCOL_ERROR_REQUEST
				readArgsFail = true
				break
			}

			cmdLen, err := parseOneCmdArgsLen(argsWithDelimiter)
			if err != nil {
				ch <- PROTOCOL_ERROR_REQUEST
				readArgsFail = true
				break
			}

			cmd := make([]byte, cmdLen+2) //这里 +2 的原因是需要一并读取 \r\n : $3\r\nset\r\n
			_, err = io.ReadFull(buf, cmd)
			if err != nil {
				ch <- PROTOCOL_ERROR_REQUEST
				readArgsFail = true
				break
			}
			cmds = append(cmds, cmd[:len(cmd)-2]) //去掉读取到  \r\n
		}

		if readArgsFail {
			continue
		}

		ch <- RedisRequet{
			Args: cmds,
		}
	}
}

/*
	解析数组的个数
		*3\r\n.....   --> 返回3
*/
func parseCmdArgsCount(header []byte) (int, error) {
	argsCountAsByte := header[1 : len(header)-2]
	argsCountAsStr := string(argsCountAsByte)
	argsCount, err := strconv.Atoi(argsCountAsStr)
	return argsCount, err
}

/*
	解析数组中的一个字符串长度
	$3\r\nset\r\n   ---> 返回3
*/
func parseOneCmdArgsLen(cmd []byte) (int, error) {
	cmdLenAsByte := cmd[1 : len(cmd)-2]
	cmdLenAsStr := string(cmdLenAsByte)
	argsCount, err := strconv.Atoi(cmdLenAsStr)
	return argsCount, err
}
