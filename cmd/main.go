package main

import (
	"os"

	goredistraning "github.com/chenjiayao/sidergo"
	"github.com/chenjiayao/sidergo/cluster"
	"github.com/chenjiayao/sidergo/config"
	"github.com/chenjiayao/sidergo/interface/server"
	"github.com/chenjiayao/sidergo/lib/logger"
	"github.com/chenjiayao/sidergo/redis"
	_ "github.com/chenjiayao/sidergo/redis/datatype"
)

func main() {
	logger.Setting()

	configFile := os.Getenv("REDIS_CONFIG")
	if configFile == "" {
		config.LoadDefaultConfig()
	} else {
		config.LoadConfig(configFile)
	}

	s := makeServer()

	goredistraning.ListenAndServe(s)
}

//TODO 根据 config 判断是否启动集群模式
func makeServer() server.Server {
	if config.Config.EnableCluster {
		return cluster.MakeCluster()
	} else {
		return redis.MakeRedisServer()
	}
}
