package packfile

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"sync"

	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/plumbing/indexfile"
)

// 32768 bytes zlib sliding window size
const ZLIB_SLIDING_WINDOW_SIZE = 32768

func (pack *Pack) Decode(packBytes []byte) error {
	reader := bytes.NewReader(packBytes)
	signature := make([]byte, 4)
	_, err := io.ReadFull(reader, signature)
	if err != nil {
		return err
	}
	pack.Signature = *(*[4]byte)(signature)

	err = binary.Read(reader, binary.BigEndian, &pack.Version)
	if err != nil {
		return err
	}

	var numEntries int32
	err = binary.Read(reader, binary.BigEndian, &numEntries)
	if err != nil {
		return err
	}
	pack.NumEntries = int(numEntries)

	pack.Objects = make([]PackObject, pack.NumEntries)
	for i := 0; i < pack.NumEntries; i++ {
		dataByte, _ := reader.ReadByte()

		var object PackObject

		// msb is a flag whether to continue read byte for size construction, 3 bits for object type and 4 bits for size
		object.Type = plumbing.ObjectType((dataByte >> 4) & 7)

		// TODO: have a threshold to prevent read large object into the memory
		object.DeflatedSize = int64(dataByte & 0x0F)
		shift := 4
		for dataByte&0x80 > 0 {
			dataByte, _ = reader.ReadByte()
			object.DeflatedSize += int64(dataByte&0x7F) << shift
			shift += 7
		}

		if object.Type == plumbing.OBJ_REF_DELTA {
			baseHash := make([]byte, 20)
			io.ReadFull(reader, baseHash)
		} else if object.Type == plumbing.OBJ_OFS_DELTA {
			getVariableLength(reader, 0, 0)
		}

		object.Size, _ = readObject(&object, reader)

		pack.Objects[i] = object
	}

	return nil
}

func (pack *Pack) DecodeWithIndex(packBytes []byte, idx *indexfile.Index) error {
	reader := bytes.NewReader(packBytes)
	signature := make([]byte, 4)
	_, err := io.ReadFull(reader, signature)
	if err != nil {
		return err
	}
	pack.Signature = *(*[4]byte)(signature)

	err = binary.Read(reader, binary.BigEndian, &pack.Version)
	if err != nil {
		return err
	}

	pack.NumEntries = len(idx.Hashes)
	pack.Objects = make([]PackObject, pack.NumEntries)

	for i, _ := range idx.Hashes  {
		offset := idx.GetOffset(i)

		reader.Seek(offset, io.SeekStart)
		dataByte, _ := reader.ReadByte()

		var object PackObject

		// msb is a flag whether to continue read byte for size construction, 3 bits for object type and 4 bits for size
		object.Type = plumbing.ObjectType((dataByte >> 4) & 7)

		// TODO: have a threshold to prevent read large object into the memory
		object.DeflatedSize = int64(dataByte & 0x0F)
		shift := 4
		for dataByte&0x80 > 0 {
			dataByte, _ = reader.ReadByte()
			object.DeflatedSize += int64(dataByte&0x7F) << shift
			shift += 7
		}

		if object.Type == plumbing.OBJ_REF_DELTA {
			baseHash := make([]byte, 20)
			io.ReadFull(reader, baseHash)
		} else if object.Type == plumbing.OBJ_OFS_DELTA {
			getVariableLength(reader, 0, 0)
		}

		object.Size, _ = readObject(&object, reader)

		pack.Objects[i] = object
	}

	return nil
}

var zlibInitBytes = []byte{0x78, 0x9c, 0x01, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0x00, 0x01}
var zlibReaderPool = sync.Pool{
	New: func() interface{} {
		r, _ := zlib.NewReader(bytes.NewReader(zlibInitBytes))
		return r
	},
}

var zlibBufferPool = sync.Pool{
	New: func() interface{} {
		bs := make([]byte, 32*1024)
		return &bs
	},
}


func readObject(writer *PackObject, reader io.Reader) (int64, error)  {
	zReader := zlibReaderPool.Get().(io.ReadCloser)
	zReader.(zlib.Resetter).Reset(reader, nil)
	defer zlibReaderPool.Put(zReader)
	defer zReader.Close()

	buffer := zlibBufferPool.Get().(*[]byte)
	written, err := io.CopyBuffer(writer, zReader, *buffer)
	zlibBufferPool.Put(buffer)

	return written, err
}

func getVariableLength(reader *bytes.Reader, length int32, shift int) int32 {
	dataByte, _ := reader.ReadByte()
	for dataByte&0x80 > 0 {
		length += (int32(dataByte&0x7F) << shift)
		shift += 7
		dataByte, _ = reader.ReadByte()
	}

	return length
}
