package unboundedchan

//TODO 没有长度的 chan 实现，用于 aof 命令写入
type UnboundedChan struct {
	In     chan<- interface{}
	Out    <-chan interface{}
	Buffer []interface{}
}
