package indexfile

import (
	"io"
	"testing"

	"github.com/liy/goe/tests"
)

func GetIndexFile(t *testing.T) io.Reader {
	file, err := tests.FS(false).Open("/topo-sort/.git/objects/pack/pack-6179faab20f2d649a12fd52aab3c8d6e32b27dcd.idx")
	if err != nil {
		t.Fatal(err)
	}

	return file
}