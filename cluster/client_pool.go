package cluster

import (
	"time"

	req "github.com/chenjiayao/sidergo/redis/request"
	"github.com/sirupsen/logrus"
)

type clientPool struct {
	ipPortPair string
	stopChan   chan struct{}
	clients    []*client
}

func MakeClientPool(ipPortPair string, num int) *clientPool {

	pool := &clientPool{
		ipPortPair: ipPortPair,
		clients:    make([]*client, num),
		stopChan:   make(chan struct{}),
	}
	for i := 0; i < num; i++ {
		pool.clients[i] = makeClient(ipPortPair)
	}

	pool.start()

	return pool
}

func (pool *clientPool) destroy() {
	pool.stopChan <- struct{}{}
}

func (pool *clientPool) start() {
	go pool.heartbeat()
}

//pool 会每隔 20 * len(clients)s 对所有的 client 进行一次 ping 请求，保证连接正常
func (pool *clientPool) heartbeat() {

	s := 20 * len(pool.clients)
	ticker := time.NewTicker(time.Duration(s) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for i := 0; i < len(pool.clients); i++ {
				client := pool.clients[i]
				if !client.isServerOnline() {
					pool.clients[i] = makeClient(client.ipPortPair)
				} else if client.IsIdle() {
					pingReq := &req.RedisRequet{
						CmdName: "ping",
						Args:    make([][]byte, 0),
					}
					r := client.SendRequestWithTimeout(pingReq, 15*time.Second)
					logrus.Info("client pool heartbeat response : ", string(r.ToContentByte()))
				}
			}
		case <-pool.stopChan:
			return
		}
	}
}
