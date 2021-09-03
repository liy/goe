package packfile

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/liy/goe/plumbing"
)

// 32768 bytes zlib sliding window size
const ZLIB_SLIDING_WINDOW_SIZE = 32768

type PackObject struct {
	Type         plumbing.ObjectType
	Data         []byte
	DeflatedSize int32
}

type Pack struct {
	Name       plumbing.Hash
	Version    int32
	Objects    []PackObject
	Signature  [4]byte
	NumEntries int
}

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
		object.DeflatedSize = int32(dataByte & 0x0F)
		shift := 4
		for dataByte&0x80 > 0 {
			dataByte, _ = reader.ReadByte()
			object.DeflatedSize += int32(dataByte&0x7F) << shift
			shift += 7
		}

		if object.Type == plumbing.OBJ_REF_DELTA {
			baseHash := make([]byte, 20)
			io.ReadFull(reader, baseHash)
		} else if object.Type == plumbing.OBJ_OFS_DELTA {
			getVariableLength(reader, 0, 0)
		}

		// object size is actually inflatted data size: deflatedObjSize >= objectSize
		// It is a good hint for zlib to read the data into the buffer without run into overflow problem
		object.Data = readObject(reader, object.DeflatedSize)

		fmt.Println(string(object.Data))

		pack.Objects[i] = object
	}

	return nil
}

func readObject(reader io.Reader, sizeHint int32) []byte {
	// TODO: re-use zlib reader
	zReader, _ := zlib.NewReader(reader)
	buffer := new(bytes.Buffer)
	chunkSize := sizeHint
	if sizeHint > ZLIB_SLIDING_WINDOW_SIZE {
		chunkSize = ZLIB_SLIDING_WINDOW_SIZE
	}
	for {
		chunk := make([]byte, chunkSize)
		// zlib has sliding window size, any data size larger than
		// the sliding window size will require multiple read.
		// The error returned is not strictly an error, but contains stream end information: EOF.
		numBytesRead, err := zReader.Read(chunk)
		buffer.Write(chunk[:numBytesRead])

		if err != nil {
			// Finish reading the compressed stream correctly
			if err.Error() == "EOF" {
				break
			} else {
				fmt.Println("Error reading zlib stream", err)
			}
		}
	}
	zReader.Close()

	return (*buffer).Bytes()
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
