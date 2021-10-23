package packfile

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/plumbing/indexfile"
)

// TODO: remove, not used

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

	for i := 0; i < int(idx.NumObjects); i++ {
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
			ReadVariableLength(reader)
		}

		decompressObjectData(object, reader)

		pack.Objects[i] = object
	}

	return nil
}

// Variable length encoding, with add 1 encoding
func ReadVariableLength(reader ByteReader) int64 {
	b, _ := reader.ReadByte()

	var v = int64(b & 0x7F)
	for b&0x80 > 0 {
		v++

		b, _ = reader.ReadByte()
		v = (v << 7) + int64(b&0x7F)
	}

	return v
}

// Variable size encoding, without 1 encoding, little endian
// This is used for reading delta deflated base size and deflated object size
func ReadVariableLengthLE(reader ByteReader) int64 {
	b, _ := reader.ReadByte()

	var v = int64(b & 0x7F)
	shift := 7
	for b&0x80 > 0 {

		b, _ = reader.ReadByte()
		v = int64(b&0x7F)<<shift + v
		shift += 7
	}

	return v
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
