package indexfile

import (
	"testing"

	"github.com/liy/goe/store"
	"github.com/liy/goe/tests"
)

func TestDecoder(t *testing.T) {
	fixture := tests.NewEmbeded(t)
	idx, err := Decode(fixture.PackIndex())
	if err != nil {
		t.Fatal(err)
	}

	tests.ToMatchSnapshot(t, idx)
}

func BenchmarkDecode(b *testing.B) {
	dotgit := store.NewDotGit("/topo-sort/.git", tests.Embeded{})
	file, _ := dotgit.PackIndex("6179faab20f2d649a12fd52aab3c8d6e32b27dcd")

	for n :=0; n <b.N; n++ {
		Decode(file)
	}
}