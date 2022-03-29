package sidergo

import (
	"fmt"
	"net"
	"os"
	"os/signal"
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

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigCh
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			server.Close()
			listener.Close()
			logrus.Info("收到关闭信号。。停止服务")
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}

		go func() {
			logrus.Info("accept new connection : ", conn.RemoteAddr())
			server.Handle(conn)
		}()
	}

}
