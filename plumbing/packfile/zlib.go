package packfile

import (
	"bytes"
	"compress/zlib"
	"io"
	"sync"
)

var zlibInitBytes = []byte{0x78, 0x9c, 0x01, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0x00, 0x01}
var zlibReaderPool = sync.Pool{
	New: func() interface{} {
		r, _ := zlib.NewReader(bytes.NewReader(zlibInitBytes))
		return r
	},
}

var zlibBufferPool = sync.Pool{
	New: func() interface{} {
		// 32768 bytes zlib sliding window size
		bs := make([]byte, 32*1024)
		return &bs
	},
}

/*
decompressObjectData deflate the zlib object data
*/
func decompressObjectData(dst io.Writer, reader io.Reader) (int64, error) {
	zReader := zlibReaderPool.Get().(io.ReadCloser)
	zReader.(zlib.Resetter).Reset(reader, nil)
	defer zlibReaderPool.Put(zReader)
	defer zReader.Close()

	buffer := zlibBufferPool.Get().(*[]byte)
	written, err := io.CopyBuffer(dst, zReader, *buffer)
	zlibBufferPool.Put(buffer)
	// written, err := io.Copy(dst, zReader)

	return written, err
}
