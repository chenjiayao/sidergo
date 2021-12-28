package main

import (
	"github.com/chenjiayao/goredistraning"
	"github.com/chenjiayao/goredistraning/interface/server"
	"github.com/chenjiayao/goredistraning/lib/logger"
)

func main() {
	logger.Setting()
	s := makeServer()
	goredistraning.ListenAndServe(s)
}

func makeServer() server.Server {
	return goredistraning.MakeRedisServer()
}
