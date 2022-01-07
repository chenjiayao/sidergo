package redis

import (
	"io"
	"os"

	"github.com/chenjiayao/goredistraning/config"
	"github.com/chenjiayao/goredistraning/interface/server"
	"github.com/chenjiayao/goredistraning/lib/logger"
	"github.com/chenjiayao/goredistraning/redis/resp"
)

// redis aof 属于写后日志，先写内存，再写日志
type AofHandler struct {
	aofChan     chan [][]byte
	redisServer server.Server
	aofFile     io.Writer
}

func (h *AofHandler) StartAof() {
	go func() {
		for cmd := range h.aofChan {
			if h.isWriteCmd(cmd[0]) {
				h.writeToAofFile(cmd)
			}
		}
	}()
}

func (h *AofHandler) writeToAofFile(cmd [][]byte) {
	asArrayResponse := resp.MakeArrayResponse(cmd)
	asBytes := asArrayResponse.ToContentByte()
	_, err := h.aofFile.Write(asBytes)
	if err != nil {
		logger.Info("write aof failed :", err.Error())
	}
}

func (h *AofHandler) LogCmd(cmd [][]byte) {
	h.aofChan <- cmd
}

func (h *AofHandler) EndAof() {
	close(h.aofChan)
}

func MakeAofHandler(server server.Server) *AofHandler {
	handler := &AofHandler{
		aofChan:     make(chan [][]byte, 4096), //TODO 这里后续应该改成使用 unboundchan 来实现无限制的 chan
		redisServer: server,
	}
	aofFileName := config.Config.AppendFilename
	file, err := os.OpenFile(aofFileName, os.O_APPEND|os.O_WRONLY, 0664)

	//TODO 这里优化 aof 文件判断
	if err != nil {
		panic(err)
	}
	handler.aofFile = file
	return handler
}

func (h *AofHandler) isWriteCmd(cmdName []byte) bool {
	return true
}
