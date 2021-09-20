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

func ToObjectType(s string) ObjectType {
	switch s {
	case "commit":
		return OBJ_COMMIT
	case "tree":
		return OBJ_TREE
	case "blob":
		return OBJ_BLOB
	case "tag":
		return OBJ_TAG
	default:
		return OBJ_INVALID
	}
}

func (t ObjectType) String() string {
	switch t {
	case OBJ_INVALID:
		return "OBJ_INVALID"
	case OBJ_COMMIT:
		return "commit"
	case OBJ_TREE:
		return "tree"
	case OBJ_BLOB:
		return "blob"
	case OBJ_TAG:
		return "tag"
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

func (h Hash) Short() string {
	return hex.EncodeToString(h[:])[:6]
}

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

func ToHash(hash string) Hash {
	bs, _ := hex.DecodeString(hash)
	return NewHash(bs)
}

func NewHash(bs []byte) Hash {
	var h Hash
	copy(h[:], bs)

	return h
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

func (o *RawObject) Write(ba []byte) (int, error) {
	o.Data = append(o.Data, ba...)
	return len(ba), nil
}

func (o RawObject) String() string {
	return fmt.Sprintf("%v\n%v %v %v", string(o.Data), o.Type, o.DeflatedSize, o.PackedSize)
}

func (o *RawObject) Hash() Hash {
	return o.hash
}

func (o *RawObject) Size() int64 {
	return o.DeflatedSize
}