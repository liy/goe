package buffer

import (
	"bufio"
	"io"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return bufio.NewReader(nil)
	},
}

func GetBuffer(reader io.Reader) *bufio.Reader {
	buf := bufferPool.Get().(*bufio.Reader)
	buf.Reset(reader)
	return buf
}

func PutBuffer(buf *bufio.Reader) {
	bufferPool.Put(buf)
}