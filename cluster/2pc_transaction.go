package cluster

import (
	"sync"
	"time"

	"github.com/chenjiayao/sidergo/interface/conn"
	"github.com/chenjiayao/sidergo/interface/request"
	"github.com/chenjiayao/sidergo/interface/response"
	redisRequest "github.com/chenjiayao/sidergo/redis/request"
)

type transaction struct {
	txID    string
	conn    conn.Conn
	cluster *Cluster
	kv      map[string]string
	wg      sync.WaitGroup
}

func (tx *transaction) begin() {
	tx.prepare()
}

//prepare 会将每个 node 的 key 上锁
// 如果上锁失败，需要执行 rollback 将所有的锁都取消
// 这个操作和具体的命令没有关系
func (tx *transaction) prepare() {

	prepareResponses := make([]response.Response, len(tx.kv))
	index := 0
	for k, _ := range tx.kv {

		key := k

		go func() {
			tx.wg.Add(1)
			ipPortPair := tx.cluster.HashRing.Hit(k)
			client := tx.cluster.PeekIdleClient(ipPortPair)

			prepareRequest := &redisRequest.RedisRequet{
				CmdName: "prepare",
				Args: [][]byte{
					[]byte(tx.txID),
					[]byte(key),
				},
			}
			r := client.SendRequestWithTimeout(prepareRequest, time.Second)
			prepareResponses[index] = r
			index++
			tx.wg.Done()
		}()
		tx.wg.Wait()
	}
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

	for k, _ := range tx.kv {
		key := k

		go func() {
			tx.wg.Add(1)
			ipPortPair := tx.cluster.HashRing.Hit(k)
			client := tx.cluster.PeekIdleClient(ipPortPair)

			//TODO这里要设计 do 命令如何传递
			commitRequest := &redisRequest.RedisRequet{
				CmdName: "commit",
				Args: [][]byte{
					[]byte(tx.txID),
					[]byte(key),
				},
			}
			r := client.SendRequestWithTimeout(commitRequest, time.Second)
			commitResponses[key] = r
			tx.wg.Done()
		}()
		tx.wg.Wait()
	}

	successKeys := make([]string, 0)
	rollback := false
	for k, r := range commitResponses {
		if r.ISOK() {
			successKeys = append(successKeys, k)
		} else {
			rollback = true
		}
	}

	if rollback {
		tx.rollbackCommit(successKeys) //undo 成功的命令
	}
	tx.unlockAllKey()

}

// 需要对每个 node 执行对应的 undo 命令，这个命令应该由调用方传递
// 命令只针对那些成功的 key,
func (tx *transaction) rollbackCommit(successKeys []string) {
	for _, key := range successKeys {

		go func(k string) {
			tx.wg.Add(1)
			ipPortPair := tx.cluster.HashRing.Hit(k)
			client := tx.cluster.PeekIdleClient(ipPortPair)

			//TODO这里要设计 undo 命令如何传递
			unlockRequest := &redisRequest.RedisRequet{
				CmdName: "rollback_commit",
				Args: [][]byte{
					[]byte(tx.txID),
					[]byte(k),
				},
			}
			client.SendRequestWithTimeout(unlockRequest, time.Second)
			tx.wg.Done()
		}(key)
	}
}

func (tx *transaction) rollbackPrepare() {
	tx.unlockAllKey()
}

//
func (tx *transaction) unlockAllKey() {

	for k, _ := range tx.kv {
		key := k // TODO 要描述清楚这里为什么要这么做
		go func() {
			tx.wg.Add(1)
			ipPortPair := tx.cluster.HashRing.Hit(k)
			client := tx.cluster.PeekIdleClient(ipPortPair)

			unlockRequest := &redisRequest.RedisRequet{
				CmdName: "transaction_unlock",
				Args: [][]byte{
					[]byte(tx.txID),
					[]byte(key),
				},
			}
			client.SendRequestWithTimeout(unlockRequest, time.Second)
			tx.wg.Done()
		}()
	}
}

func MakeTransaction(conn conn.Conn, cluster *Cluster, kv map[string]string) *transaction {
	tx := &transaction{
		conn:    conn,
		cluster: cluster,
		kv:      kv,
		wg:      sync.WaitGroup{},
	}
	return tx
}

func commit(cluster *Cluster, conn conn.Conn, clientRequest request.Request) response.Response {

	args := clientRequest.GetArgs()

	cmdName := string(args[0])
	cmdArgs := args[1:]

	command := clusterCommandRouter[cmdName]
	r := &redisRequest.RedisRequet{
		CmdName: cmdName,
		Args:    cmdArgs,
	}

	return command.CommandFunc(cluster, conn, r)
}

func undo(cluster *Cluster, conn conn.Conn, clientRequest request.Request) response.Response {
	args := clientRequest.GetArgs()

	cmdName := string(args[0])
	cmdArgs := args[1:]

	command := clusterCommandRouter[cmdName]
	r := &redisRequest.RedisRequet{
		CmdName: cmdName,
		Args:    cmdArgs,
	}
	return command.CommandFunc(cluster, conn, r)
}
