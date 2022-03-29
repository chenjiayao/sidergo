package redis

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/chenjiayao/sidergo/config"
	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/request"
	"github.com/chenjiayao/sidergo/interface/response"
	"github.com/chenjiayao/sidergo/interface/server"
	"github.com/chenjiayao/sidergo/lib/atomic"
	"github.com/chenjiayao/sidergo/lib/list"
	"github.com/chenjiayao/sidergo/parser"
	"github.com/chenjiayao/sidergo/redis/redisrequest"
	"github.com/chenjiayao/sidergo/redis/redisresponse"
	"github.com/sirupsen/logrus"
)

var _ server.Server = &RedisServer{}

// handler 实例只会有一个
type RedisServer struct {
	closed           atomic.Boolean
	rds              *RedisDBs
	aofHandler       *AofHandler
	ctx              context.Context
	cancel           context.CancelFunc
	connectedClients map[string]conn.Conn
}

///////////启动 redis 服务，
// 如果这里有 aof，那么需要加载 aof
func MakeRedisServer() *RedisServer {
	ctx, cancel := context.WithCancel(context.TODO())
	redisServer := &RedisServer{
		closed:           atomic.Boolean(0),
		ctx:              ctx,
		cancel:           cancel,
		connectedClients: make(map[string]conn.Conn),
	}

	redisServer.rds = NewDBs(redisServer)
	if config.Config.Appendonly {
		redisServer.aofHandler = MakeAofHandler(redisServer)
	}

	go redisServer.checkTimeoutConn()  //检查阻塞的链接是否超时返回 null
	go redisServer.activeExpireCycle() //定时删除过期的 key

	if config.Config.Appendonly {
		redisServer.Log()
	}

	return redisServer
}

//TODO 要处理协程退出
func (redisServer *RedisServer) activeExpireCycle() {

	/**
	1. 从过期字段中随机 20 个 key
	2. 删除这 20 个key 中过期的 key
	3. 如果过期的 key 比率超过 1/4，那么重复步骤 1
	*/
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for _, db := range redisServer.rds.DBs {
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

							if config.Config.Appendonly {
								deleteCmdRequest := &redisrequest.RedisRequet{
									CmdName: "del",
									Args: [][]byte{
										[]byte(key),
									},
								}
								redisServer.aofHandler.LogCmd(deleteCmdRequest)
							}
						}
					}

					if delKeyCount <= 5 {
						break
					}
				}
			}
		case <-redisServer.ctx.Done():
			logrus.Info("activeExpireCycle...结束")
			return
		}
	}
}

//TODO 这里要做协程退出处理，不然会导致协程泄漏
func (redisServer *RedisServer) checkTimeoutConn() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
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
							conn.SetBlockingResponse(redisresponse.NullMultiResponse)
							conn.SetBlockingExec("", nil)
							l.RemoveNode(conn) //链接已经不再阻塞，从 list 中移除
						}
						node = node.Next()
					}
					return true
				})
			}
		case <-redisServer.ctx.Done():
			logrus.Info("checkTimeoutConn结束...")
			return
		}
	}
}

func (redisServer *RedisServer) Log() {
	redisServer.aofHandler.StartAof()
}

func (redisServer *RedisServer) Handle(conn net.Conn) {

	if redisServer.closed.Get() {
		conn.Close()
		return
	}

	redisClient := MakeRedisConn(conn)
	redisServer.connectedClients[conn.RemoteAddr().String()] = redisClient

	ch := parser.ReadCommand(conn)
	//chan close 掉之后， range 直接退出
	for request := range ch {
		if request.GetErr() != nil {
			if request.GetErr() == io.EOF {
				delete(redisServer.connectedClients, redisClient.RemoteAddress())
				redisServer.closeClient(redisClient)
				return
			}

			errResponse := redisresponse.MakeErrorResponse(request.GetErr().Error())
			err := redisClient.Write(errResponse.ToErrorByte()) //返回执行命令失败，close client
			if err != nil {
				redisServer.closeClient(redisClient)
				return
			}
		}

		res := redisServer.Exec(redisClient, request)
		err := redisServer.sendResponse(redisClient, res)
		if err == io.EOF {
			break
		}
	}
}

func (redisServer *RedisServer) Exec(conn conn.Conn, request request.Request) response.Response {
	var res response.Response

	args := request.GetArgs()
	cmdName := request.GetCmdName()

	if cmdName != "auth" && !redisServer.isAuthenticated(conn) {
		res = redisresponse.MakeErrorResponse("NOAUTH Authentication required")
		return res
	}

	selectedDBIndex := conn.GetSelectedDBIndex()
	selectedDB := redisServer.rds.DBs[selectedDBIndex]

	res = selectedDB.Exec(conn, cmdName, args)

	//返回空，表示 conn 执行的是阻塞命令，当前链接被阻塞
	if res == nil {
		res = conn.GetBlockingResponse()
	}

	if res.ISOK() && config.Config.Appendonly {
		redisServer.aofHandler.LogCmd(request)
	}

	return res
}

func (redisServer *RedisServer) isAuthenticated(conn conn.Conn) bool {
	if config.Config.RequirePass == "" {
		return true
	}
	return config.Config.RequirePass == conn.GetPassword()
}

func (redisServer *RedisServer) sendResponse(redisClient conn.Conn, res response.Response) error {
	var err error
	if _, ok := res.(redisresponse.RedisErrorResponse); ok {
		err = redisClient.Write(res.ToErrorByte())
	} else {
		err = redisClient.Write(res.ToContentByte())
	}
	if err == io.EOF {
		redisServer.closeClient(redisClient)
	}
	return err
}

func (redisServer *RedisServer) LockKey(dbIndex int, key string, txID string) error {
	redisServer.rds.DBs[dbIndex].LockKey(key, txID)
	return nil
}
func (redisServer *RedisServer) UnLockKey(dbIndex int, key string, txID string) error {
	redisServer.rds.DBs[dbIndex].UnLockKey(key)
	return nil
}

// closeClient
func (redisServer *RedisServer) closeClient(client conn.Conn) {
	logrus.Info(client.RemoteAddress(), " close")
	client.Close()
}

func (redisServer *RedisServer) Close() error {
	if redisServer.closed.Get() {
		return nil
	}
	redisServer.closed.Set(true)

	for _, client := range redisServer.connectedClients {
		client.Close()
	}

	redisServer.cancel()
	if redisServer.aofHandler != nil {
		redisServer.aofHandler.EndAof()
	}
	redisServer.rds.CloseAllDB()

	time.Sleep(time.Second)
	return nil
}
