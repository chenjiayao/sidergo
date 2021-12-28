package goredistraning

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/chenjiayao/goredistraning/interface/server"
	"github.com/chenjiayao/goredistraning/lib/atomic"
	"github.com/chenjiayao/goredistraning/lib/logger"
	"github.com/chenjiayao/goredistraning/parser"
	"github.com/chenjiayao/goredistraning/redis"
)

var _ server.Server = &RedisServer{}

// handler 实例只会有一个
type RedisServer struct {
	closed atomic.Boolean
	rds    *redis.RedisDBs
}

///////////启动 redis 服务，
// 如果这里有 aof，那么需要加载 aof
func MakeRedisServer() *RedisServer {
	return &RedisServer{
		rds:    redis.NewDBs(),
		closed: atomic.Boolean(0),
	}
}

func ListenAndServe(server server.Server) {

	listener, err := net.Listen("tcp", ":8101")
	if err != nil {
		logger.Fatal("start listen failed : ", err)
	}

	logger.Info(fmt.Sprintf("start listen %s", listener.Addr().String()))
	defer func() {
		listener.Close()
		server.Close()
	}()

	var waitGroup sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("accept link")
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()
			server.Handle(conn)
		}()
	}

	//这里使用 waitGroup 的作用是：还有 conn 在处理情况下
	// 如果 redis server 关闭，那么这里需要 wait 等待已有链接处理完成。
	waitGroup.Wait()
}

func (redisServer *RedisServer) Handle(conn net.Conn) {

	if redisServer.closed.Get() {
		conn.Close()
	}

	redisClient := redis.MakeRedisConn(conn)

	ch := parser.ReadCommand(conn)
	//chan close 掉之后， range 直接退出
	for request := range ch {
		if request.Err != nil {
			if request.Err == io.EOF {
				//关闭客户端，这个
				redisServer.colseClient(redisClient)
				return
			}

			errResponse := redis.MakeErrorResponse(request.Err.Error())
			err := redisClient.Write(errResponse.ToErrorByte()) //返回执行命令失败，close client
			if err != nil {
				logger.Info("response failed: " + redisClient.RemoteAddress())
				redisServer.colseClient(redisClient)
				return
			}
		}

		logger.Info(fmt.Sprintf("get command 「%s」", request.ToStrings()))

		cmds := request.Args

		selectedDBIndex := redisClient.GetSelectedDBIndex()
		selectedDB := redisServer.rds.DBs[selectedDBIndex]
		resp := selectedDB.Exec(cmds)

		var err error
		if len(resp.ToErrorByte()) != 0 {
			err = redisClient.Write(resp.ToErrorByte())
		} else {
			err = redisClient.Write(resp.ToContentByte())
		}
		if err != nil {
			if err == io.EOF {
				redisServer.colseClient(redisClient)
			}
			continue
		}
	}
}

// colseClient
func (redisServer *RedisServer) colseClient(client *redis.RedisConn) {
	logger.Info(fmt.Sprintf("client %s closed", client.RemoteAddress()))
	client.Close()
}

func (redisServer *RedisServer) Close() error {
	logger.Info("client close....")
	return nil
}
