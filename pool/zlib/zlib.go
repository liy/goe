package zlib

import (
	"bytes"
	"compress/zlib"
	"io"
	"sync"
)

var zlibInitBytes = []byte{0x78, 0x9c, 0x01, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0x00, 0x01}
var zReaderPool = sync.Pool{
	New: func() interface{} {
		r, _ := zlib.NewReader(bytes.NewReader(zlibInitBytes))
		return r
	},
}

func GetReader(reader io.Reader) io.ReadCloser {
	zReader := zReaderPool.Get().(io.ReadCloser)
	zReader.(zlib.Resetter).Reset(reader, nil)
	return zReader
}

func PutReader(zReader io.ReadCloser) {
	zReaderPool.Put(zReader)
	zReader.Close()
}


var zlibBufferPool = sync.Pool{
	New: func() interface{} {
		// 32768 bytes zlib sliding window size
		bs := make([]byte, 32*1024)
		return &bs
	},
}

func Decompress(dst io.Writer, src io.Reader) (int64, error) {
	zReader := zReaderPool.Get().(io.ReadCloser)
	zReader.(zlib.Resetter).Reset(src, nil)
	defer zReaderPool.Put(zReader)
	defer zReader.Close()

	buffer := zlibBufferPool.Get().(*[]byte)
	written, err := io.CopyBuffer(dst, zReader, *buffer)
	zlibBufferPool.Put(buffer)
	// written, err := io.Copy(dst, zReader)

	return written, err
}