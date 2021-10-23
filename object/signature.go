package object

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Signature struct {
	Name      string
	Email     string
	TimeStamp time.Time
}

func (s *Signature) Decode(data []byte) {
	start := bytes.LastIndexByte(data, '<')
	end := bytes.LastIndexByte(data, '>')

	s.Name = string(data[:start-1])
	s.Email = string(data[start+1 : end])

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
