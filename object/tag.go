package object

import (
	"fmt"
	"strings"

	"github.com/liy/goe/errors"
	"github.com/liy/goe/plumbing"
)

type Tag struct {
	Hash       plumbing.Hash
	Target     plumbing.Hash
	TargetType plumbing.ObjectType
	Name       string
	Tagger     Signature
	Message    string
}

func NewTag(hash plumbing.Hash, message string) *Tag {
	return &Tag{
		Hash:    hash,
		Message: message,
	}
}

func (t Tag) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "object %v\n", t.Target)
	fmt.Fprintf(&sb, "type %s\n", t.TargetType)
	fmt.Fprintf(&sb, "name %s\n", t.Name)
	fmt.Fprintf(&sb, "tagger %s\n", t.Tagger)
	fmt.Fprint(&sb, "\n")
	fmt.Fprint(&sb, t.Message)

	return sb.String()
}

func DecodeTag(raw *plumbing.RawObject) (*Tag, error) {
	if raw.Type != plumbing.OBJ_TAG {
		return nil, errors.ErrRawObjectTypeWrong
	}

	t := &Tag{
		Hash: raw.Hash(),
	}

	ScanObjectData(raw.Data, func(key string, value []byte) {
		switch key {
		case "object":
			t.Target = plumbing.ToHash(string(value))
		case "type":
			t.TargetType = plumbing.ToObjectType(string(value))
		case "tag":
			t.Name = string(value)
		case "tagger":
			t.Tagger.Decode(value)
		}
	})

	return t, nil
}
