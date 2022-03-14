package redis

import (
	"io"
	"os"

	"github.com/chenjiayao/sidergo/config"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/interface/server"
	"github.com/chenjiayao/sidergo/lib/unboundedchan"
	"github.com/chenjiayao/sidergo/redis/resp"
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

	multiResponses := make([]response.Response, len(cmd))

	for i, v := range cmd {
		multiResponse := resp.MakeMultiResponse(string(v))
		multiResponses[i] = multiResponse
	}

	arrayResponse := resp.MakeArrayResponse(multiResponses)

	_, err := h.aofFile.Write(arrayResponse.ToContentByte())
	if err != nil {
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
	if os.IsNotExist(err) {
		file, err = os.Create(aofFileName)
		if err != nil {
			panic(err)
		}
	}

	//TODO 这里优化 aof 文件判断
	if err != nil {
		panic(err)
	}
	handler.aofFile = file
	return handler
}

//TODO 这里要记录所有 write 命令
func (h *AofHandler) isWriteCmd(cmdName []byte) bool {
	_, is := WriteCommands[string(cmdName)]
	return is
}
