package tests

import (
	"io"
	"testing"
)

func GetObjectFile(t *testing.T, hash string) io.Reader {
	file, err := FS(false).Open("/topo-sort/.git/objects/" + hash[:2] + "/" + hash[2:])
	if err != nil {
		t.Fatal(err)
	}

	return file
}