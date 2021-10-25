package object

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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

	buf := bufferPool.Get().(*bufio.Reader)
	buf.Reset(bytes.NewReader(raw.Data))
	defer bufferPool.Put(buf)

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
