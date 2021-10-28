package packfile

import (
	"encoding/binary"
	"io"

	"github.com/liy/goe/plumbing"
)

func Decode(file io.ReadSeeker, idx io.Reader) (*Pack, error) {
	packReader := NewPackReader(file, idx)
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

	pack.NumObjects = packReader.Index.NumObjects
	pack.Objects = make([]*plumbing.RawObject, pack.NumObjects)

	for i := 0; i < int(packReader.Index.NumObjects); i++ {
		offset := packReader.Index.GetOffsetAt(i)
		hashBytes := packReader.Index.GetHash(i)
		raw := plumbing.NewRawObject(plumbing.NewHash(hashBytes))
		packReader.readObjectAt(offset, raw)
		pack.Objects[i] = raw
	}

	return pack, nil
}