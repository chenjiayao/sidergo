package main

import (
	"os"

	"github.com/chenjiayao/goredistraning"
	"github.com/chenjiayao/goredistraning/cluster"
	"github.com/chenjiayao/goredistraning/config"
	"github.com/chenjiayao/goredistraning/interface/server"
	"github.com/chenjiayao/goredistraning/lib/logger"
	"github.com/chenjiayao/goredistraning/redis"
	_ "github.com/chenjiayao/goredistraning/redis/datatype"
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
