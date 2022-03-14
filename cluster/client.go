package cluster

import (
	"io"
	"net"
	"time"

	"github.com/chenjiayao/sidergo/interface/request"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/lib/atomic"
	req "github.com/chenjiayao/sidergo/redis/request"
	"github.com/chenjiayao/sidergo/redis/resp"
)

/*

client 的作用是，将请求发送到对应的 node ，相当于这是一个 redis-cli



cluster 会维护有一个 map[ip:port]clients


*/

type client struct {
	ipPortPair string
	conn       net.Conn
	stopChan   chan struct{}
	isIdle     atomic.Boolean
}

func makeClient(ipPortPair string) *client {

	var c *client
	n, err := net.Dial("tcp", ipPortPair)
	if err != nil {
		c = &client{
			ipPortPair: ipPortPair,
		}
	} else {
		c = &client{
			conn:       n,
			isIdle:     atomic.Boolean(1),
			ipPortPair: ipPortPair,
		}
	}
	c.Start()
	return c
}

func (c *client) SendRequest(request request.Request) chan response.Response {

	c.isIdle.Set(false)
	defer c.isIdle.Set(true)

	ch := make(chan response.Response)

	var r response.Response

	_, err := c.conn.Write(request.ToByte())
	if err != nil {
		if err == io.EOF {
			r = resp.MakeErrorResponse("server closed conn")
		} else {
			r = resp.MakeErrorResponse("unknow err")
		}
		ch <- r
		return ch
	}

	r = resp.MakeMultiResponse(c.conn.RemoteAddr().String())
	ch <- r
	return ch
}

//保持一个心跳连接，同时要判断对方是否在线
func (c *client) heartbeat() {

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if c.isServerOnline() {
				req := &req.RedisRequet{
					Args: [][]byte{
						[]byte("ping"),
					},
				}
				c.SendRequest(req)
			} else {
				conn, err := net.Dial("tcp", c.ipPortPair)
				if err != nil {
					continue
				}
				c.conn = conn
			}

		case <-c.stopChan:
			return
		}
	}
}

func (c *client) isServerOnline() bool {
	return c.conn != nil && c.conn.RemoteAddr().String() != ""
}

func (c *client) IsIdle() bool {

	return c.isServerOnline() && c.isIdle.Get()
}

func (c *client) Start() {
	go c.heartbeat()
}

func (c *client) Stop() {
	c.stopChan <- struct{}{}
}
