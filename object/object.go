package object

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/liy/goe/errors"
	"github.com/liy/goe/plumbing"
)

type ObjectDecoder interface {
	Decode(data []byte) error
}

type Object struct {
	Type string
	Size int64
	Data []byte
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return bufio.NewReader(nil)
	},
}

func ParseObjectFile(hash plumbing.Hash, repoPath string) (*Object, error) {
	// TODO: Try find it in the files system
	h := hash.String()
	p := filepath.Join(repoPath, "objects", h[:2], h[2:])
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		return nil, errors.ErrObjectNotFound
	}

	f, _ := os.Open(p)
	zReader, err := zlib.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer zReader.Close()

	buf := bufferPool.Get().(*bufio.Reader)
	buf.Reset(zReader)
	defer bufferPool.Put(buf)

	// type
	t, err := buf.ReadString(' ')
	if err != nil {
		return nil, err
	}
	
	// size
	s, err := buf.ReadString(0)
	if err != nil {
		return nil, err
	}
	size, err := strconv.ParseInt(s[:len(s)-1], 10, 64)
	if err != nil {
		return nil, err
	}

	// data
	data := make([]byte, size)
	_, err = buf.Read(data)
	if err != nil {
		return nil, err
	}

	return &Object{
		t[:len(t)-1],
		size,
		data,
	}, nil
}
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
	fmt.Fprintf(&sb, "%s <%s> %v %s", s.Name, s.Email, s.TimeStamp.Unix(), s.TimeStamp.Format("-0700"))

	return sb.String()
}

func ScanObjectData(data []byte, callback func(string, []byte)) error {
	buf := bufferPool.Get().(*bufio.Reader)
	buf.Reset(bytes.NewReader(data))
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
		callback(string(chunks[0]), chunks[1])
	}

	var msg bytes.Buffer
	for {
		line, err := buf.ReadBytes('\n')
		msg.Write(line)

		if err == io.EOF {
			break
		}
	}
	callback("message", msg.Bytes())

	return nil
}