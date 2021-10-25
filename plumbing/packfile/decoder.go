package packfile

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/plumbing/indexfile"
)

func Decode(file *os.File, idx *indexfile.Index) (*Pack, error) {
	packReader := NewPackReaderFromFile(file, idx)
	pack := new(Pack)
	
	signature := make([]byte, 4)
	_, err := io.ReadFull(packReader, signature)
	if err != nil {
		return nil, err
	}
	pack.Signature = *(*[4]byte)(signature)

	err = binary.Read(packReader, binary.BigEndian, &pack.Version)
	if err != nil {
		return nil, err
	}

	pack.NumObjects = idx.NumObjects
	pack.Objects = make([]*plumbing.RawObject, pack.NumObjects)

	for i := 0; i < int(idx.NumObjects); i++ {
		offset := idx.GetOffsetAt(i)
		hashBytes := idx.GetHash(i)
		raw := plumbing.NewRawObject(plumbing.NewHash(hashBytes))
		packReader.readObjectAt(offset, raw)
		pack.Objects[i] = raw
	}

	return pack, nil
}

// Variable length encoding, with add 1 encoding
func ReadVariableLength(reader io.ByteReader) int64 {
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
func ReadVariableLengthLE(reader io.ByteReader) int64 {
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