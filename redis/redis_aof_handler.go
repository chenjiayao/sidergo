package redis

import (
	"github.com/chenjiayao/goredistraning/helper"
	"github.com/chenjiayao/goredistraning/interface/server"
	"github.com/chenjiayao/goredistraning/lib/logger"
)

type AofHandler struct {
	aofChan     chan [][]byte
	redisServer server.Server
}

func (h *AofHandler) StartAof() {

	go func() {
		for cmd := range h.aofChan {
			logger.Info(helper.BbyteToSString(cmd))
		}
	}()
}
func (h *AofHandler) LogCmd(cmd [][]byte) {
	h.aofChan <- cmd
}

func (h *AofHandler) EndAof() {
	close(h.aofChan)
}

func MakeAofHandler(server server.Server) *AofHandler {
	return &AofHandler{
		aofChan:     make(chan [][]byte, 4096), //TODO 这里后续应该改成使用 unboundchan 来实现无限制的 chan
		redisServer: server,
	}
}
