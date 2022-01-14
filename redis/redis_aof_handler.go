package redis

import (
	"io"
	"os"

	"github.com/chenjiayao/goredistraning/config"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/interface/server"
	"github.com/chenjiayao/goredistraning/lib/logger"
	"github.com/chenjiayao/goredistraning/lib/unboundedchan"
	"github.com/chenjiayao/goredistraning/redis/resp"
)

// redis aof 属于写后日志，先写内存，再写日志
type AofHandler struct {
	aofChan     *unboundedchan.UnboundedChan
	redisServer server.Server
	aofFile     io.Writer
}

func (h *AofHandler) StartAof() {
	go func() {
		for cmd := range h.aofChan.Out {
			h.writeToAofFile(cmd)
		}
	}()
}

func (h *AofHandler) writeToAofFile(cmd [][]byte) {
	if !h.isWriteCmd(cmd[0]) {
		return
	}

	simpleResponse := resp.MakeMultiResponse(cmd)
	asArrayResponse := resp.MakeArrayResponse([]response.Response{simpleResponse})
	asBytes := asArrayResponse.ToContentByte()
	_, err := h.aofFile.Write(asBytes)
	if err != nil {
		logger.Info("write aof failed :", err.Error())
	}
}

func (h *AofHandler) LogCmd(cmd [][]byte) {
	h.aofChan.In <- cmd
}

func (h *AofHandler) EndAof() {
	defer close(h.aofChan.In)
	for cmd := range h.aofChan.Out {
		h.writeToAofFile(cmd)
	}
}

func MakeAofHandler(server server.Server) *AofHandler {
	handler := &AofHandler{
		aofChan:     unboundedchan.MakeUnboundedChan(20),
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
