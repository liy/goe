package packfile

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/plumbing/indexfile"
)

type ByteReader interface {
	ReadByte() (byte, error)
}

func (pack *Pack) Decode(packBytes []byte, idx *indexfile.Index) error {
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
	pack.Objects = make([]*plumbing.RawObject, pack.NumEntries)

	for i:=0; i<int(idx.NumObjects); i++  {
		hash := *(*[20]byte)(idx.GetHash(i))
		offset, _ := idx.GetOffset(hash)

		reader.Seek(offset, io.SeekStart)
		dataByte, _ := reader.ReadByte()

		object := plumbing.NewRawObject(hash)

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

		deflateObject(object, reader)

		pack.Objects[i] = object
	}

	return nil
}

func getVariableLength(reader ByteReader, length int32, shift int) int32 {
	dataByte, _ := reader.ReadByte()
	for dataByte&0x80 > 0 {
		length += (int32(dataByte&0x7F) << shift)
		shift += 7
		dataByte, _ = reader.ReadByte()
	}

	return length
}


// func Decode(packBytes []byte) (*Pack, error) {
// 	pack := new(Pack)

// 	reader := bytes.NewReader(packBytes)
// 	signature := make([]byte, 4)
// 	_, err := io.ReadFull(reader, signature)
// 	if err != nil {
// 		return nil, err
// 	}
// 	pack.Signature = *(*[4]byte)(signature)

// 	err = binary.Read(reader, binary.BigEndian, &pack.Version)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var numEntries int32
// 	err = binary.Read(reader, binary.BigEndian, &numEntries)
// 	if err != nil {
// 		return nil, err
// 	}
// 	pack.NumEntries = int(numEntries)

// 	pack.Objects = make([]plumbing.RawObject, pack.NumEntries)
// 	for i := 0; i < pack.NumEntries; i++ {
// 		dataByte, _ := reader.ReadByte()

// 		var object plumbing.RawObject

// 		// msb is a flag whether to continue read byte for size construction, 3 bits for object type and 4 bits for size
// 		object.Type = plumbing.ObjectType((dataByte >> 4) & 7)

// 		// TODO: have a threshold to prevent read large object into the memory
// 		object.DeflatedSize = int64(dataByte & 0x0F)
// 		shift := 4
// 		for dataByte&0x80 > 0 {
// 			dataByte, _ = reader.ReadByte()
// 			object.DeflatedSize += int64(dataByte&0x7F) << shift
// 			shift += 7
// 		}

// 		if object.Type == plumbing.OBJ_REF_DELTA {
// 			baseHash := make([]byte, 20)
// 			io.ReadFull(reader, baseHash)
// 		} else if object.Type == plumbing.OBJ_OFS_DELTA {
// 			getVariableLength(reader, 0, 0)
// 		}

// 		deflateObject(&object, reader)

// 		pack.Objects[i] = object
// 	}

// 	return pack, nil
// }