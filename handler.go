package goredistraning

import (
	"fmt"
	"io"
	"net"

	"github.com/chenjiayao/goredistraning/interface/handler"
	"github.com/chenjiayao/goredistraning/lib/atomic"
	"github.com/chenjiayao/goredistraning/lib/logger"
	"github.com/chenjiayao/goredistraning/parser"
	"github.com/chenjiayao/goredistraning/redis"
)

var _ handler.Handler = RedisHandler{}

// handler 实例只会有一个
type RedisHandler struct {
	closed atomic.Boolean
}

func (rh RedisHandler) Handle(conn net.Conn) {

	if rh.closed.Get() {
		conn.Close()
	}

	redisClient := redis.MakeRedisConn(conn)

	ch := parser.ReadCommand(conn)
	//chan close 掉之后， range 直接退出
	for request := range ch {
		if request.Err != nil {
			if request.Err == io.EOF {
				//关闭客户端，这个
				rh.colseClient(redisClient)
				return
			}

			errResponse := redis.MakeErrorResponse(request.Err.Error())
			err := redisClient.Write(errResponse.ToErrorByte()) //返回执行命令失败，close client
			if err != nil {
				logger.Info("response failed: " + redisClient.RemoteAddress())
				rh.colseClient(redisClient)
				return
			}
		}

		logger.Info(fmt.Sprintf("get command 「%s」", request.ToStrings()))

		cmds := request.Args
		//TODO 执行命令

		sr := redis.MakeSimpleResponse(cmds)
		redisClient.Write(sr.ToContentByte())
	}

}

// colseClient
func (rh RedisHandler) colseClient(client *redis.RedisConn) {
	logger.Info(fmt.Sprintf("client %s closed", client.RemoteAddress()))
	client.Close()
}

func (rh RedisHandler) Close() error {
	logger.Info("client close....")
	return nil
}
