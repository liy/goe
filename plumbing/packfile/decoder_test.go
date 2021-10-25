package packfile

import (
	"testing"

	"github.com/liy/goe/fixtures"
	"github.com/liy/goe/plumbing/indexfile"
	"github.com/liy/goe/utils"
)

func TestDecode(t *testing.T) {
	fixture := fixtures.NewRepositoryFixture("topo-sort")
	idxFile := fixture.IndexFile("../../repos")
	packFile := fixture.PackFile("../../repos")
	
	idx, _ := indexfile.Decode(idxFile)
	pack, err := Decode(packFile, idx)
	if err != nil {
		t.Fatal(err)
	}

	utils.ToMatchSnapshot(t, pack)
}