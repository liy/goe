package indexfile

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/liy/goe/utils"
)

// func Decode(indexBytes []byte) (*Index, error) {
func Decode(reader io.Reader) (*Index, error) {
	idx := new(Index)
	// reader := bytes.NewReader(indexBytes)

	// Check magic header
	magicBytes := make([]byte, 4)
	_, err := io.ReadFull(reader, magicBytes)
	utils.CheckIfError(err)
	if string(magicBytes) != "\377tOc" {
		return nil, fmt.Errorf("invalid IDX header, only version 2 supported: %q", string(magicBytes))
	}

	// version
	err = binary.Read(reader, binary.BigEndian, &idx.Version)
	if err != nil {
		return nil, err
	}

	// fanout and bucket
	err = binary.Read(reader, binary.BigEndian, &idx.Fanout[0])
	if err != nil {
		return nil, err
	}
	idx.Buckets[0] = idx.Fanout[0]
	for i := 1; i < 256; i++ {
		err = binary.Read(reader, binary.BigEndian, &idx.Fanout[i])
		if err != nil {
			return nil, err
		}
		idx.Buckets[i] = idx.Fanout[i] - idx.Fanout[i-1]
	}

	idx.NumObjects = idx.Fanout[255]
	// hashes
	idx.Hashes = make([]byte, idx.NumObjects*20)
	_, err = io.ReadFull(reader, idx.Hashes)
	if err != nil {
		return nil, err
	}

	// crc
	idx.CRC = make([]byte, idx.NumObjects*4)
	_, err = io.ReadFull(reader, idx.CRC)
	if err != nil {
		return nil, err
	}

	// offsets 32 
	idx.Offset32 = make([]byte, idx.NumObjects*4)
	_, err = io.ReadFull(reader, idx.Offset32)
	if err != nil {
		return nil, err
	}
	
	// Check if any offset64, msb is 1
	var numOffset64 int
	for i := 0; i < int(idx.NumObjects); i+=4 {
		if idx.Offset32[i]&128 > 0 {
			numOffset64++
		}
	}

	// offset 64
	if numOffset64 > 0 {
		idx.Offset64 = make([]byte, numOffset64*8)
		_, err = io.ReadFull(reader, idx.Offset64)
		if err != nil {
			return nil, err
		}
	}

	// For reverse hash lookup, using offset
	idx.ReverseHash = make(map[uint64]uint32, idx.NumObjects)
	for i := 0; i < int(idx.NumObjects); i++ {
		offset := idx.getOffset(i)
		idx.ReverseHash[uint64(offset)] = uint32(i)
	}
	
	return idx, nil
}
