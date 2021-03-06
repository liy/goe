package object

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/liy/goe/errors"
	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/pool/buffer"
)

/*
Tag represents an annotated tag in git.
*/
type Tag struct {
	// The hash of the annotated tag object
	Hash       plumbing.Hash
	// The object that tag points to
	Target     plumbing.Hash
	// The object's type that tag points to
	TargetType plumbing.ObjectType
	// Name of the tag
	Name       string
	// The creator of the tag
	Tagger     Signature
	// Message of the tag
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

	buf := buffer.GetBuffer(bytes.NewReader(raw.Data))
	defer buffer.PutBuffer(buf)


	for {
		line, err := buf.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		// message starts from the first empty line
		line = bytes.TrimRight(line, "\n")
		if len(line) == 0 {
			break
		}

		chunks := bytes.SplitN(line, []byte{' '}, 2)

		switch string(chunks[0]) {
		case "object":
			t.Target = plumbing.ToHash(string(chunks[1]))
		case "type":
			t.TargetType = plumbing.ToObjectType(string(chunks[1]))
		case "tag":
			t.Name = string(chunks[1])
		case "tagger":
			t.Tagger.Decode(chunks[1])
		}
	}

	return t, nil
}
