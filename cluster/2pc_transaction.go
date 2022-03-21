package cluster

import (
	"sync"
	"time"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/request"
	"github.com/chenjiayao/sidergo/interface/response"
	redisRequest "github.com/chenjiayao/sidergo/redis/request"
	"github.com/chenjiayao/sidergo/redis/resp"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

func init() {
	RegisterClusterExecCommand("prepare", ExecPrepare, nil)
	RegisterClusterExecCommand("commit", ExecCommit, nil)
	RegisterClusterExecCommand("undo", ExecUndo, nil)
	RegisterClusterExecCommand("transaction_unlock", ExecTransactionUnlock, nil)
}

//redis 中只有 mset 需要自用分布式事务处理
type transaction struct {
	txID           string
	conn           conn.Conn
	cluster        *Cluster
	undoRequests   []request.Request
	commitRequests []request.Request
	wg             sync.WaitGroup
}

func (tx *transaction) begin() {
	tx.prepare()
}

//prepare 会将每个 node 的 key 上锁
// 如果上锁失败，需要执行 rollback 将所有的锁都取消
// 这个操作和具体的命令没有关系
func (tx *transaction) prepare() {

	prepareResponses := make([]response.Response, 0)
	for _, re := range tx.commitRequests {
		tx.wg.Add(1)

		key := re.GetKey()
		prepareRequest := &redisRequest.RedisRequet{
			CmdName: "prepare",
			Args: [][]byte{
				[]byte(key),
				[]byte(tx.txID),
			},
		}

		ipPortPair := tx.cluster.HashRing.Hit(key)

		if tx.cluster.Self.IsSelf(ipPortPair) {
			prepareResponses = append(prepareResponses, ExecPrepare(tx.cluster, tx.conn, prepareRequest))
			tx.wg.Done()
		} else {
			go func(key string) {
				client := tx.cluster.PeekIdleClient(ipPortPair)
				prepareResponses = append(prepareResponses, client.SendRequestWithTimeout(prepareRequest, time.Second))
				tx.wg.Done()
			}(key)
		}
	}
	tx.wg.Wait()

	logrus.Info("prepare wait ok")
	for _, r := range prepareResponses {
		if !r.ISOK() {
			tx.rollbackPrepare()
		}
	}
	tx.commit()
}

//需要对每个 node 执行对应的命令，这个命令应该由调用方传递
// 如果命令执行失败，那么要对成功的那部分执行 undo 命令，
//这里 commit 不管成功失败都不可以应该去掉每个 node 的 lock key
func (tx *transaction) commit() {
	commitResponses := make(map[string]response.Response)

	for _, commitRequest := range tx.commitRequests {
		tx.wg.Add(1)

		key := commitRequest.GetKey()

		ipPortPair := tx.cluster.HashRing.Hit(commitRequest.GetKey())

		if tx.cluster.Self.IsSelf(ipPortPair) {
			commitResponses[key] = ExecCommit(tx.cluster, tx.conn, commitRequest)
			tx.wg.Done()
		} else {
			go func(key string) {
				client := tx.cluster.PeekIdleClient(ipPortPair)
				commitResponses[key] = client.SendRequestWithTimeout(commitRequest, time.Second)
				tx.wg.Done()
			}(key)
		}
	}
	tx.wg.Wait()

	successKeys := make(map[string]string)
	rollback := false
	for k, r := range commitResponses {
		if r.ISOK() {
			successKeys[k] = ""
		} else {
			rollback = true
		}
	}

	if rollback {
		tx.rollbackCommit(successKeys) //undo 成功的命令
	}
	logrus.Info("commit ok")
	tx.unlockAllKey()

}

// 需要对每个 node 执行对应的 undo 命令，这个命令应该由调用方传递
// 命令只针对那些成功的 key,
func (tx *transaction) rollbackCommit(successKeys map[string]string) {

	for _, undoRequest := range tx.undoRequests {

		key := undoRequest.GetKey()

		_, exist := successKeys[key]
		if !exist {
			continue
		}

		tx.wg.Add(1)
		ipPortPair := tx.cluster.HashRing.Hit(key)

		if tx.cluster.Self.IsSelf(ipPortPair) {
			ExecUndo(tx.cluster, tx.conn, undoRequest)
			tx.wg.Done()
		} else {
			go func() {
				client := tx.cluster.PeekIdleClient(ipPortPair)
				client.SendRequestWithTimeout(undoRequest, time.Second)
				tx.wg.Done()
			}()
		}
	}

	tx.wg.Wait()
}

func (tx *transaction) rollbackPrepare() {
	tx.unlockAllKey()
}

//
func (tx *transaction) unlockAllKey() {

	for _, undoRequest := range tx.undoRequests {
		tx.wg.Add(1)

		k := undoRequest.GetKey()
		ipPortPair := tx.cluster.HashRing.Hit(k)
		unlockRequest := &redisRequest.RedisRequet{
			CmdName: "transaction_unlock",
			Args: [][]byte{
				[]byte("transaction_unlock"),
				[]byte(tx.txID),
				[]byte(k),
			},
		}

		if tx.cluster.Self.IsSelf(ipPortPair) {
			ExecTransactionUnlock(tx.cluster, tx.conn, unlockRequest)
			tx.wg.Done()
		} else {
			go func() {
				client := tx.cluster.PeekIdleClient(ipPortPair)
				client.SendRequestWithTimeout(unlockRequest, time.Second)
				tx.wg.Done()
			}()
		}
	}
	tx.wg.Wait()
	logrus.Info("unlock done")
}

func (tx *transaction) generateUniqueID() string {
	return xid.New().String()
}

func MakeTransaction(conn conn.Conn, cluster *Cluster, undoRequests []request.Request, commitRequests []request.Request) *transaction {
	tx := &transaction{
		conn:           conn,
		cluster:        cluster,
		undoRequests:   undoRequests,
		commitRequests: commitRequests,
		wg:             sync.WaitGroup{},
	}
	tx.txID = tx.generateUniqueID()
	return tx
}

///////////另一个node 转发请求到这里执行
func ExecCommit(cluster *Cluster, conn conn.Conn, clientRequest request.Request) response.Response {

	args := clientRequest.GetArgs()
	cmdName := string(args[0])

	command := clusterCommandRouter[cmdName]
	cmdRequest := &redisRequest.RedisRequet{
		CmdName: cmdName,
		Args:    args[1:],
	}
	logrus.Info(command.CmdName)
	return command.CommandFunc(cluster, conn, cmdRequest)
}

//执行 undo 操作，注意不取消 unlock
func ExecUndo(cluster *Cluster, conn conn.Conn, clientRequest request.Request) response.Response {
	args := clientRequest.GetArgs()
	undoCommandName := string(args[0])

	undoCommand := clusterCommandRouter[undoCommandName]
	cmdRequest := &redisRequest.RedisRequet{
		CmdName: undoCommandName,
		Args:    args[1:],
	}
	return undoCommand.CommandFunc(cluster, conn, cmdRequest)
}

//txid 和 key
func ExecPrepare(cluster *Cluster, conn conn.Conn, clientRequest request.Request) response.Response {

	args := clientRequest.GetArgs()
	key := string(args[0])
	txID := string(args[1])
	selectedDBIndex := conn.GetSelectedDBIndex()

	//锁定key
	err := cluster.Self.RedisServer.LockKey(selectedDBIndex, key, txID)
	if err != nil {
		return resp.MakeErrorResponse(err.Error())
	}
	return resp.OKSimpleResponse
}

func ExecTransactionUnlock(cluster *Cluster, conn conn.Conn, clientRequest request.Request) response.Response {
	args := clientRequest.GetArgs()
	txID := string(args[0])
	key := string(args[1])
	selectedDBIndex := conn.GetSelectedDBIndex()
	cluster.Self.RedisServer.UnLockKey(selectedDBIndex, key, txID)
	return resp.OKSimpleResponse
}
