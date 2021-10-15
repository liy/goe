package object

import (
	"fmt"
	"strings"

	"github.com/liy/goe/plumbing"
)

type Tag struct {
	Hash plumbing.Hash
	Target plumbing.Hash
	TargetType plumbing.ObjectType
	Name string
	Tagger Signature
	Message string
}

func NewTag(hash plumbing.Hash, message string) *Tag{
	return &Tag {
		Hash: hash,
		Message: message,
	}
}

func (t *Tag) Decode(data []byte) error {
	return ScanObjectData(data, func(key string, value []byte) {
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
	t := &Tag{
		Hash: raw.Hash(),
	}
	err := t.Decode(raw.Data)
	return t, err
}