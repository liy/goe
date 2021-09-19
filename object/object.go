package object

import (
	"bufio"
	"sync"

	"github.com/liy/goe/plumbing"
)

type Object interface {
	Decode(raw *plumbing.RawObject) error
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return bufio.NewReader(nil)
	},
}
