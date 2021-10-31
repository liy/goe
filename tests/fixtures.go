package tests

import (
	"io"
	"net/http"
	"path"
	"testing"

	"github.com/liy/goe/store"
)

type Fixture struct {
	Url string
	packHashes []string

	DotGit *store.DotGit
	t *testing.T
}

func NewEmbeded(t *testing.T) *Fixture {
	return &Fixture{
		"https://gitlab.com/liyss/goe-fixtures/-/raw/main/topo-sort.tar.gz",
		[]string{"6179faab20f2d649a12fd52aab3c8d6e32b27dcd"},
		store.NewDotGit("/topo-sort/.git", Embeded{}),
		t,
	}
}

func (f *Fixture) PackIndex() http.File {
	file, err := f.DotGit.PackIndex(f.packHashes[0])
	if err != nil {
		f.t.Fatal(err)
	}
	return file
}

func (f *Fixture) Pack() http.File {
	file, err := f.DotGit.Pack(f.packHashes[0])
	if err != nil {
		f.t.Fatal(err)
	}
	return file
}

func (f *Fixture) GetObjectFile(hash string) io.Reader {
	file, err := f.DotGit.Open(path.Join("objects", hash[:2], hash[2:]))
	if err != nil {
		f.t.Fatal(err)
	}

	return file
}