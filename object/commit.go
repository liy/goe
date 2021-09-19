package object

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/liy/goe/plumbing"
)

type Signature struct {
	Name string
	Email string
	TimeStamp time.Time
}

func (s *Signature) Decode(data []byte) {
	start := bytes.LastIndexByte(data, '<')
	end := bytes.LastIndexByte(data, '>')

	s.Name = string(data[:start-1])
	s.Email = string(data[start+1:end])

	// parse date time
	chunks := bytes.Split(data[end+2:], []byte{' '})
	ts, err := strconv.ParseInt(string(chunks[0]), 10, 64)
	if err != nil {
		return
	}

	// Timezone not used?
	// tz, err := strconv.Atoi(string(chunks[1]))
	// if err != nil {
	// 	return
	// }

	s.TimeStamp = time.Unix(ts, 0)
}

func (s Signature) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s <%s> %s", s.Name, s.Email, s.TimeStamp)

	return sb.String()
}

type Commit struct {
	Hash plumbing.Hash
	Tree plumbing.Hash
	Parents []plumbing.Hash
	Author Signature
	Commiter Signature
	Message string
}

func (c *Commit) Decode(raw *plumbing.RawObject) (err error) {
	if raw.Type != plumbing.OBJ_COMMIT {
		return fmt.Errorf("not commit object")
	}

	r := bytes.NewReader(raw.Data)
	reader := bufferPool.Get().(*bufio.Reader)
	defer bufferPool.Put(r)
	reader.Reset(r)

	for {
		line, err := reader.ReadBytes('\n')
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
			c.Commiter.Decode(chunks[1])
		}
	}

	var msg bytes.Buffer
	for {
		line, err := reader.ReadBytes('\n')
		msg.Write(line)

		if err == io.EOF {
			break
		}
	}
	c.Message = msg.String()

	return nil
}

func (c Commit) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "tree %v\n", c.Tree)
	for _, p := range c.Parents {
		fmt.Fprintf(&sb, "parent %v\n", p)
	}
	fmt.Fprintf(&sb, "author %s\n", c.Author)
	fmt.Fprintf(&sb, "commiter %s\n", c.Commiter)
	fmt.Fprint(&sb, "\n")
	fmt.Fprint(&sb, c.Message)

	return sb.String()
}

func DecodeCommit(raw *plumbing.RawObject) (*Commit, error) {
	c := &Commit{
		Hash: raw.Hash(),
	}
	err := c.Decode(raw)
	return c, err
}