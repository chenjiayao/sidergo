package goredistraning

import (
	"fmt"
	"net"
	"sync"

	"github.com/chenjiayao/sidergo/config"
	"github.com/chenjiayao/sidergo/interface/server"
)

func ListenAndServe(server server.Server) {

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Config.Bind, config.Config.Port))
	if err != nil {
	}

	if config.Config.Appendonly {
		server.Log()
	}

	defer func() {
		listener.Close()
		server.Close()
	}()

	var waitGroup sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()
			server.Handle(conn)
		}()
	}

	//这里使用 waitGroup 的作用是：还有 conn 在处理情况下
	// 如果 redis server 关闭，那么这里需要 wait 等待已有链接处理完成。
	waitGroup.Wait()
}
