package indexfile

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/utils"
)

type Index struct {
	Hashes   []plumbing.Hash
	CRC      [][4]byte
	Offset32 [][4]byte
	Offset64 []byte
}

const MSB_MASK_32 = uint32(1) << 31

func (idx *Index) Decode(indexBytes []byte) error {
	reader := bytes.NewReader(indexBytes)

	// Check magic header
	magicBytes := make([]byte, 4)
	_, err := io.ReadFull(reader, magicBytes)
	utils.CheckIfError(err)
	if string(magicBytes) != "\377tOc" {
		return fmt.Errorf("invalid IDX header, only version 2 supported: %q", string(magicBytes))
	}

	// Fanout, number of objects
	// Just need total object number for now
	reader.Seek(256*4, io.SeekCurrent)
	var numObjects uint32
	binary.Read(reader, binary.BigEndian, &numObjects)

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

func (idx *Index) GetOffset(position int) uint64 {
	// Mask out the MSB to construct offset
	offset := binary.BigEndian.Uint32(idx.Offset32[position][:]) & ^MSB_MASK_32
	// If the MSB is set, the rest 31bit offset is actually the index to the Offset64 bytes
	if idx.Offset32[position][0]&128 > 0 {
		return binary.BigEndian.Uint64(idx.Offset64[offset : offset+8])
	} else {
		return uint64(offset)
	}
}
