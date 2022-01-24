package redis

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/chenjiayao/goredistraning/config"
	"github.com/chenjiayao/goredistraning/interface/conn"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/interface/server"
	"github.com/chenjiayao/goredistraning/lib/atomic"
	"github.com/chenjiayao/goredistraning/lib/list"
	"github.com/chenjiayao/goredistraning/lib/logger"
	"github.com/chenjiayao/goredistraning/parser"
	"github.com/chenjiayao/goredistraning/redis/resp"
)

var _ server.Server = &RedisServer{}

// handler 实例只会有一个
type RedisServer struct {
	closed     atomic.Boolean
	rds        *RedisDBs
	aofHandler *AofHandler
}

///////////启动 redis 服务，
// 如果这里有 aof，那么需要加载 aof
func MakeRedisServer() *RedisServer {
	redisServer := &RedisServer{
		closed: atomic.Boolean(0),
	}

	redisServer.rds = NewDBs(redisServer)
	if config.Config.Appendonly {
		redisServer.aofHandler = MakeAofHandler(redisServer)
	}

	go redisServer.CheckTimeoutConn()

	return redisServer
}

func (redisServer *RedisServer) CheckTimeoutConn() {
	for {
		for _, db := range redisServer.rds.DBs {
			db.BlockingKeys.Range(func(key, value interface{}) bool {
				l, _ := value.(*list.List)
				node := l.First()
				for {
					if node == nil {
						break
					}
					element := node.Element()
					conn, _ := element.(conn.Conn)
					blockAt := conn.GetBlockAt()
					blockTime := conn.GetMaxBlockTime()
					if blockTime == 0 {
						continue
					}

					if time.Until(blockAt) > time.Duration(blockTime) {
						conn.SetBlockingResponse(resp.NullMultiResponse)
					}
					node = node.Next()
				}
				return true
			})
		}
	}

}

func (redisServer *RedisServer) Log() {
	redisServer.aofHandler.StartAof()
}

func (redisServer *RedisServer) Handle(conn net.Conn) {

	if redisServer.closed.Get() {
		conn.Close()
	}

	redisClient := MakeRedisConn(conn)

	ch := parser.ReadCommand(conn)
	//chan close 掉之后， range 直接退出
	for request := range ch {
		if request.Err != nil {
			if request.Err == io.EOF {
				redisServer.closeClient(redisClient)
				return
			}

			errResponse := resp.MakeErrorResponse(request.Err.Error())
			err := redisClient.Write(errResponse.ToErrorByte()) //返回执行命令失败，close client
			if err != nil {
				logger.Info("response failed: " + redisClient.RemoteAddress())
				redisServer.closeClient(redisClient)
				return
			}
		}

		var res response.Response
		var err error

		cmd := request.Args
		cmdName := redisServer.parseCommand(request.Args)
		args := cmd[1:]

		if cmdName != "auth" && !redisServer.isAuthenticated(redisClient) {
			res = resp.MakeErrorResponse("NOAUTH Authentication required")
			err := redisServer.sendResponse(redisClient, res)
			if err == io.EOF {
				break
			}
			continue
		}

		selectedDBIndex := redisClient.GetSelectedDBIndex()
		selectedDB := redisServer.rds.DBs[selectedDBIndex]

		res = selectedDB.Exec(redisClient, cmdName, args)

		//返回空，表示 conn 执行的是阻塞命令，当前链接被阻塞
		if res == nil {
			res = redisClient.GetBlockingResponse()
		}

		err = redisServer.sendResponse(redisClient, res)
		if res.ISOK() && config.Config.Appendonly {
			redisServer.aofHandler.LogCmd(request.Args)
		}

		if err == io.EOF {
			break
		}
	}
}

func (redisServer *RedisServer) isAuthenticated(redisClient *RedisConn) bool {
	return config.Config.RequirePass == redisClient.GetPassword()
}

func (redisServer *RedisServer) sendResponse(redisClient *RedisConn, res response.Response) error {
	var err error
	if _, ok := res.(resp.RedisErrorResponse); ok {
		err = redisClient.Write(res.ToErrorByte())
	} else {
		err = redisClient.Write(res.ToContentByte())
	}
	if err == io.EOF {
		redisServer.closeClient(redisClient)
	}
	return err
}

//从请求数据中解析出 redis 命令
func (redisServer *RedisServer) parseCommand(cmd [][]byte) string {
	cmdName := string(cmd[0])
	return strings.ToLower(cmdName)
}

// closeClient
func (redisServer *RedisServer) closeClient(client *RedisConn) {
	logger.Info(fmt.Sprintf("client %s closed", client.RemoteAddress()))
	client.Close()
}

func (redisServer *RedisServer) Close() error {
	logger.Info("server close....")
	redisServer.aofHandler.EndAof()
	return nil
}
