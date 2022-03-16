package cluster

//分布式事务 2pc 实现
// 用于解决 mset 命令
/*

1. 向所有的 key 发出 prepare 请求
2. 向所有的 key 发送 commit 请求


*/
type transaction struct {
}

func (tx *transaction) prepare() {

}

func (tx *transaction) commit() {

}

func (tx *transaction) rollback() {

}
