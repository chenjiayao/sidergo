package cluster

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/chenjiayao/sidergo/interface/request"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/lib/atomic"
	"github.com/chenjiayao/sidergo/redis/resp"
	"github.com/sirupsen/logrus"
)

/*

client 的作用是，将请求发送到对应的 node ，相当于这是一个 redis-cli



cluster 会维护有一个 map[ip:port]clients

*/

type client struct {
	ipPortPair string
	conn       net.Conn
	isIdle     atomic.Boolean //1/true是空闲，0/false是忙碌

	stopChan chan struct{}
}

func makeClient(ipPortPair string) *client {

	var c *client
	n, err := net.Dial("tcp", ipPortPair)
	if err != nil {
		logrus.Info("client to node server failed : ", err, ipPortPair)
		c = &client{
			ipPortPair: ipPortPair,
			conn:       nil,
			isIdle:     atomic.Boolean(1),
			stopChan:   make(chan struct{}),
		}
	} else {
		c = &client{
			conn:       n,
			isIdle:     atomic.Boolean(1),
			ipPortPair: ipPortPair,
			stopChan:   make(chan struct{}),
		}
	}
	return c
}

func (c *client) SendRequestWithTimeout(request request.Request, timeout time.Duration) response.Response {

	c.isIdle.Set(false)
	defer c.isIdle.Set(true)

	var r response.Response

	logrus.Info("send command : ", string(request.ToByte()))
	_, err := c.conn.Write(request.ToByte())
	if err != nil {
		if err == io.EOF {
			r = resp.MakeErrorResponse("server closed conn")
		} else {
			r = resp.MakeErrorResponse(err.Error())
		}
		return r
	}

	b, err := c.parse(c.conn)
	if err != nil {
		if err == io.EOF {
			c.conn.Close()
			return resp.MakeErrorResponse("server closed")
		}
		return resp.MakeErrorResponse(err.Error())
	}
	r = resp.MakeReidsRawByteResponse(b)
	return r
}

//只要解析出 []byte
func (c *client) parse(reader io.Reader) ([]byte, error) {

	typ := make([]byte, 1)
	reader.Read(typ)

	var err error
	var resp []byte

	t := string(typ)
	switch t {
	case "$": // 多行字符串
		resp, err = c.parseMulti(reader)
	case ":": // 数字
		resp, err = c.parseNumber(reader)
	case "+": //单行字符串
		resp, err = c.parseLine(reader)
	case "*": //数组
		resp, err = c.parseArray(reader)
	case "-": //错误
		resp, err = c.parseError(reader)
	default:
		resp = nil
		err = errors.New("protocol err")
	}
	return resp, err
}

func (c *client) parseMulti(reader io.Reader) ([]byte, error) {
	buf := bufio.NewReader(reader)
	b, err := buf.ReadBytes('\n') // 3\r\n
	if err != nil {
		return nil, err
	}

	bs := string(b[:len(b)-2]) //去掉 \r\n
	l, err := strconv.ParseInt(bs, 10, 64)
	if err != nil {
		return nil, errors.New("protocol err")
	}

	if l == -1 {
		return []byte("$-1\r\n"), nil
	}

	ret := make([]byte, 0)

	ret = append(ret, []byte("$")...)
	ret = append(ret, b...)

	s := make([]byte, l+2)
	io.ReadFull(buf, s)
	ret = append(ret, s...)
	return ret, nil
}

//"3\r\n$3\r\nget\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"
func (c *client) parseArray(reader io.Reader) ([]byte, error) {
	buf := bufio.NewReader(reader)

	b, err := buf.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	var resp []byte
	resp = append(resp, []byte("*")...)
	resp = append(resp, b...)

	l, err := strconv.ParseInt(string(b[:len(b)-2]), 10, 64)
	if err != nil {
		return nil, errors.New("protocol err")
	}

	for i := int64(1); i <= l; i++ {
		b, err := c.parse(buf)
		if err != nil {
			return nil, err
		}
		logrus.Info(string(b))
		resp = append(resp, b...)
	}
	return resp, nil
}

func (c *client) parseLine(io io.Reader) ([]byte, error) {
	ret := make([]byte, 0)

	ret = append(ret, []byte("+")...)
	buf := bufio.NewReader(io)
	b, err := buf.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	ret = append(ret, b...)
	return ret, nil
}

func (c *client) parseError(io io.Reader) ([]byte, error) {
	ret := make([]byte, 0)

	ret = append(ret, []byte("-")...)
	buf := bufio.NewReader(io)
	b, err := buf.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	ret = append(ret, b...)
	return ret, nil
}
func (c *client) parseNumber(io io.Reader) ([]byte, error) {
	ret := make([]byte, 0)

	ret = append(ret, []byte(":")...)

	buf := bufio.NewReader(io)
	b, err := buf.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	ret = append(ret, b...)
	return ret, nil
}

func (c *client) isServerOnline() bool {
	return c.conn != nil
}

func (c *client) IsIdle() bool {
	return c.isServerOnline() && c.isIdle.Get()
}

func (c *client) Stop() {
	c.stopChan <- struct{}{}
}
