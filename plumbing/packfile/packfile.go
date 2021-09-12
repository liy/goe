package packfile

import (
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

func (pr *PackReader) ReadObjectAt(offset int64, raw *plumbing.RawObject)  error {
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

	// TODO: un-delta the delta raw 
	if raw.Type == plumbing.OBJ_REF_DELTA {
		baseHash := make([]byte, 20)
		io.ReadFull(pr, baseHash)
	} else if raw.Type == plumbing.OBJ_OFS_DELTA {
		getVariableLength(pr, 0, 0)
	}

	_, err := deflateObject(raw, pr)

	return err
}

func (pr *PackReader) GetObject(hash plumbing.Hash) (*plumbing.RawObject, error) {
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