package packfile

import (
	"bytes"
	"fmt"
	"io"

	"github.com/liy/goe/plumbing"
)

type PackObject struct {
	Type         plumbing.ObjectType
	Data         []byte
	DeflatedSize int64
	Size         int64
}

func (o *PackObject) Write(ba []byte) (int, error) {
	o.Data = append(o.Data, ba...)
	return len(ba), nil
}

func (o *PackObject) GetTypeName() string {
	switch o.Type {
	case plumbing.OBJ_INVALID:
		return "OBJ_INVALID"
	case plumbing.OBJ_COMMIT:
		return "OBJ_COMMIT"
	case plumbing.OBJ_TREE:
		return "OBJ_TREE"
	case plumbing.OBJ_BLOB:
		return "OBJ_BLOB"
	case plumbing.OBJ_TAG:
		return "OBJ_TAG"
	// 5 is reserved for future expansion
	case plumbing.OBJ_OFS_DELTA:
		return "OBJ_OFS_DELTA"
	case plumbing.OBJ_REF_DELTA:
		return "OBJ_REF_DELTA"
	default:
		return "error type"
	}
}

func (o PackObject) String() string {
	if o.Type < 5 {
		return fmt.Sprintf("%v %v %v\n%v\n", o.GetTypeName(), int(o.DeflatedSize), int(o.Size), string(o.Data))
	} else {
		return fmt.Sprintf("%v %v %v\n", o.GetTypeName(), int(o.DeflatedSize), int(o.Size))
	}
}

type Pack struct {
	Name       plumbing.Hash
	Version    int32
	Objects    []PackObject
	Signature  [4]byte
	NumEntries int
}

func (pack *Pack) TestReadObjectAt(offset int64, packBytes []byte) PackObject {
	reader := bytes.NewReader(packBytes)
	reader.Seek(offset, io.SeekStart)
	dataByte, _ := reader.ReadByte()

	var object PackObject

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
		getVariableLength(reader, 0, 0)
	}

	object.Size, _ = readObject(&object, reader)

	return object
}