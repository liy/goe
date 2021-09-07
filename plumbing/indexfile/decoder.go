package indexfile

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/utils"
)

func (idx *Index) Decode(indexBytes []byte) error {
	reader := bytes.NewReader(indexBytes)

	// Check magic header
	magicBytes := make([]byte, 4)
	_, err := io.ReadFull(reader, magicBytes)
	utils.CheckIfError(err)
	if string(magicBytes) != "\377tOc" {
		return fmt.Errorf("invalid IDX header, only version 2 supported: %q", string(magicBytes))
	}

	// version
	binary.Read(reader, binary.BigEndian, &idx.Version)

	// fanout and bucket
	binary.Read(reader, binary.BigEndian, &idx.Fanout[0])
	idx.Buckets[0] = idx.Fanout[0]
	for i := 1; i < 256; i++ {
		binary.Read(reader, binary.BigEndian, &idx.Fanout[i])
		idx.Buckets[i] = idx.Fanout[i] - idx.Fanout[i-1]
	}

	// reader.Seek(256*4 , io.SeekCurrent)
	// var numObjects uint32
	// binary.Read(reader, binary.BigEndian, &numObjects)

	numObjects := idx.Fanout[255]
	// hashes
	idx.Hashes = make([]plumbing.Hash, numObjects)
	for i := 0; i < int(numObjects); i++ {
		bytes := make([]byte, 20)
		_, err := io.ReadFull(reader, bytes)
		utils.CheckIfError(err)
		idx.Hashes[i] = *(*[20]byte)(bytes)
	}

	// crc
	idx.CRC = make([][4]byte, numObjects)
	for i := 0; i < int(numObjects); i++ {
		bytes := make([]byte, 4)
		_, err := io.ReadFull(reader, bytes)
		utils.CheckIfError(err)
		idx.CRC[i] = *(*[4]byte)(bytes)
	}

	// offsets, 32 and 64
	var numOffset64 int
	idx.Offset32 = make([][4]byte, numObjects)
	for i := 0; i < int(numObjects); i++ {
		bytes := make([]byte, 4)
		_, err := io.ReadFull(reader, bytes)
		utils.CheckIfError(err)
		idx.Offset32[i] = *(*[4]byte)(bytes)

		if bytes[0]&128 > 0 {
			numOffset64++
		}
	}
	// offset 64 if any
	if numOffset64 > 0 {
		idx.Offset64 = make([]byte, numOffset64*8)
		io.ReadFull(reader, idx.Offset64)
	}

	return nil
}
