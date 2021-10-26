package plumbing

import (
	"bytes"
	"encoding/hex"
	"io"
	"testing"

	"github.com/liy/goe/tests"
	"github.com/stretchr/testify/assert"
)

func TestToHash(t *testing.T) {
	expected := [20]byte{77, 98, 175, 51, 208, 42, 30, 35, 149, 40, 34, 95, 213, 148, 57, 230, 53, 201, 250, 234}
	result := ([20]byte)(ToHash("4d62af33d02a1e239528225fd59439e635c9faea"))
	if !bytes.Equal(expected[:], result[:]) {
		t.Fatalf("ToHash mismatch")
	}
}

func TestBytes(t *testing.T) {
	hash := [20]byte{77, 98, 175, 51, 208, 42, 30, 35, 149, 40, 34, 95, 213, 148, 57, 230, 53, 201, 250, 234}
	hashStr := "4d62af33d02a1e239528225fd59439e635c9faea"
	expected, _ := hex.DecodeString(hashStr)

	if !bytes.Equal(expected, hash[:]) {
		t.Fatalf("ToHash mismatch")
	}
}

func TestRawObjectWrite(t *testing.T) {
	data := []byte{12, 22, 33, 44}
	byteReader := bytes.NewBuffer(data)

	raw := new(RawObject)
	raw.Data = []byte{2}
	numBytes, err := io.Copy(raw, byteReader)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []byte{2, 12, 22, 33, 44}, raw.Data, "Can use RawObject as a writer")
	assert.Equal(t, 4, int(numBytes), "Correct number bytes written to raw object")
}

func TestReadFile(t *testing.T) {
	raw := GetRawObject(t, "f2010ee942a47bec0ca7e8f04240968ea5200735")
	tests.ToMatchSnapshot(t, raw)
}