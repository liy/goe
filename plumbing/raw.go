package plumbing

import (
	"encoding/hex"
	"fmt"
	"io"
	"strconv"

	"github.com/liy/goe/pool/buffer"
	"github.com/liy/goe/pool/zlib"
)

type ObjectType int8

const (
	OBJ_INVALID ObjectType = 0
	OBJ_COMMIT  ObjectType = 1
	OBJ_TREE    ObjectType = 2
	OBJ_BLOB    ObjectType = 3
	OBJ_TAG     ObjectType = 4

	// 5 is reserved for future expansion

	// These 2 object types only exist in pack file
	// They represents two different delta object type:
	// offset delta and reference delta
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

var ZeroHash Hash

func (h Hash) Bytes() *[20]byte {
	return (*[20]byte)(&h)
}

func (h Hash) Slice() []byte {
	return (*h.Bytes())[:]
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
	hash Hash
	// RawType can be normal OBJ_COMMIT, OBJ_TAG, OBJ_BLOB or OBJ_TREE objects
	// but it can also be a delta of another object which means it can als be
	// OBJ_OFS_DELTA and OBJ_REF_DELTA
	RawType ObjectType
	// Type can only be normal git object: OBJ_COMMIT, OBJ_TAG, OBJ_BLOB or OBJ_TREE objects
	Type         ObjectType
	Data         []byte
	// DeflatedSize is the size of the object after deflating
	DeflatedSize int64
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
	return fmt.Sprintf("%v\n%v %v", string(o.Data), o.Type, o.DeflatedSize)
}

func (o *RawObject) Hash() Hash {
	return o.hash
}

/*
Size is the size of the object after deflating
*/
func (o *RawObject) Size() int64 {
	return o.DeflatedSize
}

func (raw *RawObject) LooseRead(reader io.Reader) error {
	zReader := zlib.GetReader(reader)
	defer zlib.PutReader(zReader)

	buf := buffer.GetBuffer(zReader)
	defer buffer.PutBuffer(buf)

	// type
	t, err := buf.ReadString(' ')
	if err != nil {
		return err
	}
	raw.Type = ToObjectType(string(t[:len(t)-1]))
	
	// Raw type is as same as the object type, since it is not loaded from pack file
	raw.RawType = raw.Type

	// size
	s, err := buf.ReadString(0)
	if err != nil {
		return err
	}
	raw.DeflatedSize, err = strconv.ParseInt(s[:len(s)-1], 10, 64)
	if err != nil {
		return err
	}

	// data
	raw.Data = make([]byte, raw.DeflatedSize)
	_, err = buf.Read(raw.Data)

	return err
}
