package main

import (
	"os"

	"github.com/chenjiayao/goredistraning"
	"github.com/chenjiayao/goredistraning/config"
	"github.com/chenjiayao/goredistraning/interface/server"
	"github.com/chenjiayao/goredistraning/lib/logger"
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
	return goredistraning.MakeRedisServer()
}
