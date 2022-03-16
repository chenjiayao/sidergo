package main

import (
	"os"

	"github.com/chenjiayao/sidergo"
	"github.com/chenjiayao/sidergo/cluster"
	"github.com/chenjiayao/sidergo/config"
	"github.com/chenjiayao/sidergo/interface/server"
	"github.com/sirupsen/logrus"

	"github.com/chenjiayao/sidergo/redis"
	_ "github.com/chenjiayao/sidergo/redis/datatype"
)

func main() {

	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetReportCaller(true)

	configFile := os.Getenv("REDIS_CONFIG")
	if configFile == "" {
		config.LoadDefaultConfig()
	} else {
		config.LoadConfig(configFile)
	}

	s := makeServer()

	sidergo.ListenAndServe(s)
}

//TODO 根据 config 判断是否启动集群模式
func makeServer() server.Server {
	if config.Config.EnableCluster {
		return cluster.MakeCluster()
	} else {
		return redis.MakeRedisServer()
	}
}
