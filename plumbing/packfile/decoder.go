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