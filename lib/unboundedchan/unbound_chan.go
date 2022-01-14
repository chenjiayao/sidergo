package unboundedchan

//TODO 没有长度的 chan 实现，用于 aof 命令写入
type UnboundedChan struct {
	In     chan<- [][]byte
	Out    <-chan [][]byte
	Buffer [][][]byte
}

func MakeUnboundedChan(initial int) *UnboundedChan {

	in := make(chan [][]byte, initial)
	out := make(chan [][]byte, initial)
	buffer := make([][][]byte, initial)

	go func() {
		defer close(out)
	loop:
		for {
			packet, ok := <-in
			if !ok {
				break loop
			}
			select {
			case out <- packet:
				continue
			default:
			}

			buffer = append(buffer, packet)

			for len(buffer) > 0 {
				select {
				case packet, ok := <-in:
					if !ok {
						break loop
					}
					buffer = append(buffer, packet)
				case out <- buffer[0]:
					buffer = buffer[1:]
				}
			}
		}

		for len(buffer) > 0 {
			out <- buffer[0]
			buffer = buffer[1:]
		}
	}()

	return &UnboundedChan{
		In:     in,
		Out:    out,
		Buffer: buffer,
	}
}
