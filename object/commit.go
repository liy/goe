package object

import (
	"fmt"
	"strings"

	"github.com/liy/goe/plumbing"
)


type Commit struct {
	Hash plumbing.Hash
	Tree plumbing.Hash
	Parents []plumbing.Hash
	Author Signature
	Committer Signature
	Message string
}

func (c *Commit) Decode(data []byte) error {
	return ScanObjectData(data, func(key string, value []byte) {
		switch key {
		case "tree":
			c.Tree = plumbing.ToHash(string(value))
		case "parent":
			c.Parents = append(c.Parents, plumbing.ToHash(string(value)))
		case "author":
			c.Author.Decode(value)
		case "committer":
			c.Committer.Decode(value)
		case "message":
			c.Message = string(value)
		}
	})
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
	c := &Commit{
		Hash: raw.Hash(),
	}
	err := c.Decode(raw.Data)
	return c, err
}