package packfile

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/plumbing/indexfile"
	"github.com/liy/goe/utils"
)

type Pack struct {
	Name       plumbing.Hash
	Version    int32
	Objects    []*plumbing.RawObject
	Signature  [4]byte
	NumEntries int
}

type PackReader struct {
	*indexfile.Index
	file io.ReadSeeker
	path string
	cache utils.Cache
	// one byte
	s []byte
}

func NewPackReader(packPath string) *PackReader {
	file, _ := os.Open(packPath)
	return &PackReader{
		Index: indexfile.NewIndex(packPath[:len(packPath)-4] + "idx"),
		file: file,
		path: packPath,
		cache: utils.NewLRU(int64(1 * 1024 * 1024)),
		s: make([]byte, 1),
	}
}

func (pr *PackReader) ReadByte() (byte, error) {
	_, err := pr.file.Read(pr.s)
	return pr.s[0], err
}

func (pr *PackReader) Read(b []byte) (int, error) {
	return pr.file.Read(b)
}

func (pr *PackReader) Seek(offset int64, whence int) (int64, error) {
	return pr.file.Seek(offset, whence)
}

func (pr *PackReader) ReadObjectAt(offset int64, raw *plumbing.RawObject) error {
	pr.Seek(offset, io.SeekStart)
	dataByte, _ := pr.ReadByte()

	// msb is a flag whether to continue read byte for size construction, 3 bits for raw type and 4 bits for size
	raw.Type = plumbing.ObjectType((dataByte >> 4) & 7)

	// TODO: have a threshold to prevent read large raw into the memory
	raw.DeflatedSize = int64(dataByte & 0x0F)
	shift := 4
	for dataByte&0x80 > 0 {
		dataByte, _ = pr.ReadByte()
		raw.DeflatedSize += int64(dataByte&0x7F) << shift
		shift += 7
	}
	raw.PackedSize = raw.DeflatedSize

	if raw.Type == plumbing.OBJ_REF_DELTA {
		hb := make([]byte, 20)
		io.ReadFull(pr, hb)
		baseHash := plumbing.NewHash(hb)

		// Traverse to read base object
		rawBase, err := pr.ReadObject(baseHash)
		if err != nil {
			return err
		}
		
		buffer := bytes.NewBuffer(make([]byte, 0))
		decompressObjectData(buffer, pr)

		// For error checking only
		baseSize := ReadVariableLengthLE(buffer)
		if rawBase.DeflatedSize != baseSize {
			return fmt.Errorf("wrong base object size")
		}
		raw.DeflatedSize = ReadVariableLengthLE(buffer)

		baseReader := bytes.NewReader(rawBase.Data)

		err = pr.DeltaPatch(buffer, baseReader, raw)
		if err != nil {
			return err
		}
	} else if raw.Type == plumbing.OBJ_OFS_DELTA {
		// Read the negative offset to the base object 
		baseOffset := offset - ReadVariableLength(pr)

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
			err := pr.ReadObjectAt(baseOffset, rawBase)
			if err != nil {
				return nil
			}
			pr.cache.Add(rawBase)
			pr.Seek(resume, io.SeekStart)
		}

		buffer := bytes.NewBuffer(make([]byte, 0))
		decompressObjectData(buffer, pr)

		// For error checking only
		baseSize := ReadVariableLengthLE(buffer)
		if rawBase.DeflatedSize != baseSize {
			return fmt.Errorf("wrong base object size")
		}
		raw.DeflatedSize = ReadVariableLengthLE(buffer)

		baseReader := bytes.NewReader(rawBase.Data)
		err := pr.DeltaPatch(buffer, baseReader, raw)
		if err != nil {
			return err
		}
	} else {
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
	if !ok  {
		return nil, fmt.Errorf("cannot find %s raw object in pack: %s",  hash.Short(), pr.path)
	}

	raw := plumbing.NewRawObject(hash)
	err := pr.ReadObjectAt(offset, raw)
	if err != nil {
		return nil, err
	}
	pr.cache.Add(raw)

	return raw, nil
}

func (pr *PackReader) DeltaPatch(deltaReader *bytes.Buffer, baseReader *bytes.Reader, dest io.Writer) error {	
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
				size = uint32(b) << 8 | size
			}
			if (cmdByte & 0x40) != 0 {
				b, _ := deltaReader.ReadByte()
				size = uint32(b) << 16 | size
			}
			if size == 0 {
				size = 0x10000
			}

			baseReader.Seek(int64(offset), io.SeekStart)
			io.CopyN(dest, baseReader, int64(size))
		} else if (cmdByte & 0x80) == 0 && cmdByte != 0 { // copy from data after command byte
			size := uint(cmdByte)
			// Read out the size of the data to be inserted, the data is followed
			io.CopyN(dest, deltaReader, int64(size))
		} else { // end of delta
			break;
		}
	}

	return nil
}