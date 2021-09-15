package plumbing

import (
	"encoding/hex"
	"fmt"
)

type ObjectType int8

const (
	OBJ_INVALID ObjectType = 0
	OBJ_COMMIT  ObjectType = 1
	OBJ_TREE    ObjectType = 2
	OBJ_BLOB    ObjectType = 3
	OBJ_TAG     ObjectType = 4
	// 5 is reserved for future expansion
	OBJ_OFS_DELTA ObjectType = 6
	OBJ_REF_DELTA ObjectType = 7
)

func (t ObjectType) String() string {
	switch t {
	case OBJ_INVALID:
		return "OBJ_INVALID"
	case OBJ_COMMIT:
		return "OBJ_COMMIT"
	case OBJ_TREE:
		return "OBJ_TREE"
	case OBJ_BLOB:
		return "OBJ_BLOB"
	case OBJ_TAG:
		return "OBJ_TAG"
	// 5 is reserved for future expansion
	case OBJ_OFS_DELTA:
		return "OBJ_OFS_DELTA"
	case OBJ_REF_DELTA:
		return "OBJ_REF_DELTA"
	default:
		return "OBJ_UNKNOWN"
	}
}

type Hash [20]byte

func (h *Hash) Bytes() *[20]byte {
	return (*[20]byte)(h)
}

func (h *Hash) Short() string {
	return hex.EncodeToString(h.Bytes()[:])[:6]
}

func ToHash(hash string) Hash {
	bs, _ := hex.DecodeString(hash)
	return NewHash(bs)
}

func NewHash(bs []byte) Hash {
	return *(*[20]byte)(bs)
}

type RawObject struct {
	hash 		 Hash
	Type         ObjectType
	Data         []byte
	DeflatedSize int64
	// Size of the object in the pack file.
	// It might be different from the real "DeflatedSize"
	PackedSize int64
}

func NewRawObject(hash Hash) *RawObject {
	return &RawObject{
		hash: hash,
	}
}

// Used by zlib deflate process for writing object 
func (o *RawObject) Write(ba []byte) (int, error) {
	o.Data = append(o.Data, ba...)
	return len(ba), nil
}

func (o RawObject) String() string {
	if o.Type < 5 {
		return fmt.Sprintf("%v %v\n%v\n", o.Type, int(o.DeflatedSize), string(o.Data))
	} else {
		return fmt.Sprintf("%v %v\n", o.Type, int(o.DeflatedSize))
	}
}

func (o *RawObject) Hash() Hash {
	return o.hash
}

func (o *RawObject) Size() int64 {
	return o.DeflatedSize
}