package packfile

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/liy/goe/errors"
	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/plumbing/indexfile"
	"github.com/liy/goe/utils"
)

type Pack struct {
	Name       plumbing.Hash
	Version    int32
	Objects    []*plumbing.RawObject
	Signature  [4]byte
	NumObjects uint32
}

type PackReader struct {
	*indexfile.Index
	file      io.ReadSeeker
	bufReader *bufio.Reader
	cache     utils.Cache
	offset    int64
}

func NewPackReader(packReader io.ReadSeeker, idxReader io.Reader) *PackReader {
	return &PackReader{
		Index:     indexfile.NewIndex(idxReader),
		file:      packReader,
		bufReader: bufio.NewReader(packReader),
		cache:     utils.NewLRU(int64(5 * 1024 * 1024)),
	}
}

func (pr *PackReader) ReadByte() (byte, error) {
	b, err := pr.bufReader.ReadByte()
	if err == nil {
		pr.offset++
	}
	return b, err
}

func (pr *PackReader) Read(b []byte) (int, error) {
	n, err := pr.bufReader.Read(b)
	if err == nil {
		pr.offset += int64(n)
	}
	return n, err
}

func (pr *PackReader) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekCurrent && offset == 0 {
		return pr.offset, nil
	}

	var err error
	pr.offset, err = pr.file.Seek(offset, whence)
	pr.bufReader.Reset(pr.file)

	return pr.offset, err
}

func (pr *PackReader) readObjectAt(offset int64, raw *plumbing.RawObject) error {
	pr.Seek(offset, io.SeekStart)
	dataByte, _ := pr.ReadByte()

	// msb is a flag whether to continue read byte for size construction, 3 bits for raw type and 4 bits for size
	raw.RawType = plumbing.ObjectType((dataByte >> 4) & 7)

	// TODO: have a threshold to prevent read large raw into the memory
	raw.DeflatedSize = int64(dataByte & 0x0F)
	shift := 4
	for dataByte&0x80 > 0 {
		dataByte, _ = pr.ReadByte()
		raw.DeflatedSize += int64(dataByte&0x7F) << shift
		shift += 7
	}
	// raw.PackedSize = raw.DeflatedSize

	// base object is specified by hash
	if raw.RawType == plumbing.OBJ_REF_DELTA {
		hb := make([]byte, 20)
		pr.Read(hb)
		baseHash := plumbing.NewHash(hb)

		// Traverse to read base object
		rawBase, err := pr.ReadObject(baseHash)
		if err != nil {
			return err
		}
		// Raw object type is actually its base type
		raw.Type = rawBase.Type

		buffer := bytes.NewBuffer(make([]byte, 0))
		decompressObjectData(buffer, pr)

		// For error checking only
		baseSize := utils.ReadVariableSize(buffer)
		if rawBase.DeflatedSize != baseSize {
			return fmt.Errorf("wrong base object size")
		}
		raw.DeflatedSize = utils.ReadVariableSize(buffer)

		baseReader := bytes.NewReader(rawBase.Data)
		err = pr.deltaPatch(buffer, baseReader, raw)
		if err != nil {
			return err
		}
	} else if raw.RawType == plumbing.OBJ_OFS_DELTA { // base object is specified by offset
		// Read the negative offset to the base object
		baseOffset := offset - utils.ReadVariableOffset(pr)

		hb, ok := pr.GetHashFromOffset(uint64(baseOffset))
		if !ok {
			return nil
		}
		baseHash := plumbing.NewHash(hb)

		// Traverse to read base object
		var rawBase *plumbing.RawObject
		if item, ok := pr.cache.Get(baseHash); ok {
			rawBase = item.(*plumbing.RawObject)
		} else {
			rawBase = plumbing.NewRawObject(baseHash)
			resume, _ := pr.Seek(0, io.SeekCurrent)
			err := pr.readObjectAt(baseOffset, rawBase)
			if err != nil {
				return nil
			}
			pr.cache.Add(rawBase)
			pr.Seek(resume, io.SeekStart)
		}
		// Raw object type is actually its base type
		raw.Type = rawBase.Type

		buffer := bytes.NewBuffer(make([]byte, 0))
		decompressObjectData(buffer, pr)

		// For error checking only
		baseSize := utils.ReadVariableSize(buffer)
		if rawBase.DeflatedSize != baseSize {
			return fmt.Errorf("wrong base object size")
		}
		raw.DeflatedSize = utils.ReadVariableSize(buffer)

		baseReader := bytes.NewReader(rawBase.Data)
		err := pr.deltaPatch(buffer, baseReader, raw)
		if err != nil {
			return err
		}
	} else {
		raw.Type = raw.RawType
		decompressObjectData(raw, pr)
	}

	return nil
}

func (pr *PackReader) ReadObject(hash plumbing.Hash) (*plumbing.RawObject, error) {
	item, ok := pr.cache.Get(hash)
	if ok {
		return item.(*plumbing.RawObject), nil
	}

	offset, ok := pr.Index.GetOffset(hash)
	if !ok {
		return nil, errors.ErrObjectNotFound
	}

	raw := plumbing.NewRawObject(hash)
	err := pr.readObjectAt(offset, raw)
	if err != nil {
		return nil, err
	}
	pr.cache.Add(raw)

	return raw, nil
}

func (pr *PackReader) deltaPatch(deltaReader *bytes.Buffer, baseReader *bytes.Reader, dest io.Writer) error {
	// Reconstruct the object data from base object
	for {
		cmdByte, _ := deltaReader.ReadByte()
		// copy from base object
		if (cmdByte & 0x80) != 0 {
			// decode offset
			var offset uint32
			if (cmdByte & 0x01) != 0 {
				b, _ := deltaReader.ReadByte()
				offset = uint32(b)
			}
			if (cmdByte & 0x02) != 0 {
				b, _ := deltaReader.ReadByte()
				offset = (uint32(b) << 8) | offset
			}
			if (cmdByte & 0x04) != 0 {
				b, _ := deltaReader.ReadByte()
				offset = (uint32(b) << 16) | offset
			}
			if (cmdByte & 0x08) != 0 {
				b, _ := deltaReader.ReadByte()
				offset = (uint32(b) << 24) | offset
			}

			// decode size
			var size uint32
			if (cmdByte & 0x10) != 0 {
				b, _ := deltaReader.ReadByte()
				size = uint32(b)
			}
			if (cmdByte & 0x20) != 0 {
				b, _ := deltaReader.ReadByte()
				size = uint32(b)<<8 | size
			}
			if (cmdByte & 0x40) != 0 {
				b, _ := deltaReader.ReadByte()
				size = uint32(b)<<16 | size
			}
			if size == 0 {
				size = 0x10000
			}

			baseReader.Seek(int64(offset), io.SeekStart)
			io.CopyN(dest, baseReader, int64(size))
		} else if (cmdByte&0x80) == 0 && cmdByte != 0 { // copy from data after command byte
			size := uint(cmdByte)
			// Read out the size of the data to be inserted, the data is followed
			io.CopyN(dest, deltaReader, int64(size))
		} else { // end of delta
			break
		}
	}

	return nil
}
