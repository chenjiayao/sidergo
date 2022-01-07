package main

import (
	"os"

	"github.com/chenjiayao/goredistraning"
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

func makeServer() server.Server {
	return redis.MakeRedisServer()
}
