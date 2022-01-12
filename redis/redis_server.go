package redis

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/chenjiayao/goredistraning/config"
	"github.com/chenjiayao/goredistraning/interface/response"
	"github.com/chenjiayao/goredistraning/interface/server"
	"github.com/chenjiayao/goredistraning/lib/atomic"
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

	redisServer.rds = NewDBs()
	if config.Config.Appendonly {
		redisServer.aofHandler = MakeAofHandler(redisServer)
	}
	return redisServer
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
		err = redisServer.sendResponse(redisClient, res)

		//这里进行判断，如果是修改命令，并且执行成功，那么需要看看 key 是否被 watch
		if res.ISOK() && !redisClient.IsInMultiState() && redisServer.isWriteCommand(cmdName) {
			//检查 key 是否被 watch
			if selectedDB.isWatched(cmdName) {
				redisClient.DirtyCAS(true)
			}
		}

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

func (redisServer *RedisServer) isWriteCommand(cmdName string) bool {
	_, is := WriteCommands[cmdName]
	return is
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
