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

	go redisServer.checkTimeoutConn()  //检查阻塞的链接是否超时返回 null
	go redisServer.activeExpireCycle() //定时删除过期的 key

	return redisServer
}

//TODO 要处理协程退出
func (redisServe *RedisServer) activeExpireCycle() {

	/**
	1. 从过期字段中随机 20 个 key
	2. 删除这 20 个key 中过期的 key
	3. 如果过期的 key 比率超过 1/4，那么重复步骤 1
	*/
	for {
		for _, db := range redisServe.rds.DBs {

			for {
				delKeyCount := 0
				for i := 0; i < 20; i++ {
					k := db.TtlMap.RandomKey()
					if k == nil {
						break
					}
					key := k.(string)
					v, _ := db.TtlMap.Get(key)
					ttlUnixTimestamp := v.(int64)

					if time.Now().Unix() > ttlUnixTimestamp {
						//删除这个key
						db.Dataset.Del(key)
						delKeyCount++
					}
				}

				if delKeyCount <= 5 {
					break
				}
			}
		}
	}
}

//TODO 这里要做协程退出处理，不然会导致协程泄漏
func (redisServer *RedisServer) checkTimeoutConn() {
	for {
		for _, db := range redisServer.rds.DBs {
			db.BlockingKeys.Range(func(key, value interface{}) bool {
				l, _ := value.(*list.List)
				node := l.HeadNode()

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

					// time.Now().Sub(blockAt) --> time - blockAt
					// time.Until(blockAt) --> blockAt.Sub(time.Now()) --> blockAt - time.Now()
					//  blockTime < time.Now() - blockAt
					if time.Since(blockAt).Seconds() > float64(blockTime) {
						conn.SetBlockingResponse(resp.NullMultiResponse)
						conn.SetBlockingExec("", nil)
						l.RemoveNode(conn) //链接已经不再阻塞，从 list 中移除
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

		if len(request.Args) == 0 {
			continue
		}

		cmd := request.Args
		cmdName := redisServer.parseCommand(cmd)
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
