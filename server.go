package sidergo

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/chenjiayao/sidergo/config"
	"github.com/chenjiayao/sidergo/interface/server"
	"github.com/sirupsen/logrus"
)

func ListenAndServe(server server.Server) {

	address := fmt.Sprintf("%s:%d", config.Config.Bind, config.Config.Port)
	logrus.Info("listen at:", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logrus.Fatal("start server failed ", err)
	}

	if config.Config.Appendonly {
		server.Log()
	}

	defer server.Close()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigCh
		logrus.Info("收到关闭信号。。停止服务")
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			listener.Close()
		}
	}()

	var waitGroup sync.WaitGroup

	for {
		conn, err := listener.Accept()
		logrus.Info("accept new connection : ", conn.RemoteAddr())
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
