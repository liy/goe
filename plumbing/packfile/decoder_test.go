package packfile

import (
	"io"
	"testing"

	"github.com/liy/goe/store"
	"github.com/liy/goe/tests"
)

func TestDecode(t *testing.T) {
	fixture := tests.NewEmbeded(t)
	pack, err := Decode(fixture.Pack(), fixture.PackIndex())
	if err != nil {
		t.Fatal(err)
	}

	tests.ToMatchSnapshot(t, pack)
}

func BenchmarkDecode(b *testing.B) {
	dotgit := store.NewDotGit("/topo-sort/.git", tests.Embeded{})
	pack, _ := dotgit.Pack("6179faab20f2d649a12fd52aab3c8d6e32b27dcd")
	packIndex, _ := dotgit.PackIndex("6179faab20f2d649a12fd52aab3c8d6e32b27dcd")

	for n := 0; n < b.N; n++ {
		// Reset readers
		pack.Seek(0, io.SeekStart)
		packIndex.Seek(0, io.SeekStart)

		Decode(pack, packIndex)
	}
}