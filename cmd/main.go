package main

import (
	"github.com/chenjiayao/goredistraning"
	"github.com/chenjiayao/goredistraning/interface/handler"
	"github.com/chenjiayao/goredistraning/lib/logger"
)

func main() {

	logger.Setting()
	goredistraning.ListenAndServe(makeHandler())
}

func makeHandler() handler.Handler {
	return goredistraning.RedisHandler{}
}
