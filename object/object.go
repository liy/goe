package object

import (
	"bufio"
	"sync"
)


var bufferPool = sync.Pool{
	New: func() interface{} {
		return bufio.NewReader(nil)
	},
}