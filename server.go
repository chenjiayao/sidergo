package goredistraning

import (
	"net"
	"sync"

	"github.com/chenjiayao/goredistraning/interface/handler"
	"github.com/chenjiayao/goredistraning/lib/logger"
)

func ListenAndServe(handler handler.Handler) {

	listener, err := net.Listen("tcp", ":8101")
	if err != nil {
		logger.Fatal("start listen failed : ", err)
	}

	defer func() {
		listener.Close()
		handler.Close()
	}()

	var waitGroup sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("accept link")
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()
			handler.Handle(conn)
		}()
	}

	//这里使用 waitGroup 的作用是：还有 conn 在处理情况下
	// 如果 redis server 关闭，那么这里需要 wait 等待已有链接处理完成。
	waitGroup.Wait()
}
