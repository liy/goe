package packfile

import (
	"testing"

	"github.com/liy/goe/plumbing/indexfile"
	"github.com/liy/goe/tests"
)

func TestDecode(t *testing.T) {
	fixture := tests.GetFixture("topo-sort")
	idxFile := fixture.IndexFile("../../repos")
	packFile := fixture.PackFile("../../repos")
	
	idx, _ := indexfile.Decode(idxFile)
	pack, err := Decode(packFile, idx)
	if err != nil {
		t.Fatal(err)
	}

	tests.ToMatchSnapshot(t, pack)
}