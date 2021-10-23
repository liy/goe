package object

import (
	"bufio"
	"bytes"
	"io"
	"sync"
)

type ObjectDecoder interface {
	Decode(data []byte) error
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return bufio.NewReader(nil)
	},
}

func ScanObjectData(data []byte, callback func(string, []byte)) {
	buf := bufferPool.Get().(*bufio.Reader)
	buf.Reset(bytes.NewReader(data))
	defer bufferPool.Put(buf)

	for {
		line, err := buf.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		// message starts from the first empty line
		line = bytes.TrimRight(line, "\n")
		if len(line) == 0 {
			break
		}

		chunks := bytes.SplitN(line, []byte{' '}, 2)
		callback(string(chunks[0]), chunks[1])
	}

	var msg bytes.Buffer
	for {
		line, err := buf.ReadBytes('\n')
		msg.Write(line)

		if err == io.EOF {
			break
		}
	}
	callback("message", msg.Bytes())
}
