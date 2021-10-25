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

type Commit struct {
	Hash      plumbing.Hash
	Tree      plumbing.Hash
	Parents   []plumbing.Hash
	Author    Signature
	Committer Signature
	Message   string
}

func (c *Commit) GetCompareValue() int {
	return int(c.Author.TimeStamp.Unix())
}

func (c Commit) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "tree %v\n", c.Tree)
	for _, p := range c.Parents {
		fmt.Fprintf(&sb, "parent %v\n", p)
	}
	fmt.Fprintf(&sb, "author %s\n", c.Author)
	fmt.Fprintf(&sb, "committer %s\n", c.Committer)
	fmt.Fprint(&sb, "\n")
	fmt.Fprint(&sb, c.Message)

	return sb.String()
}

func DecodeCommit(raw *plumbing.RawObject) (*Commit, error) {
	if raw.Type != plumbing.OBJ_COMMIT {
		return nil, errors.ErrRawObjectTypeWrong
	}

	c := &Commit{
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
			case "tree":
				c.Tree = plumbing.ToHash(string(chunks[1]))
			case "parent":
				c.Parents = append(c.Parents, plumbing.ToHash(string(chunks[1])))
			case "author":
				c.Author.Decode(chunks[1])
			case "committer":
				c.Committer.Decode(chunks[1])
		}
	}

	var msg bytes.Buffer
	for {
		line, err := buf.ReadBytes('\n')
		msg.Write(line)

		if err == io.EOF {
			break
		}
	}
	c.Message = msg.String()

	return c, nil
}
