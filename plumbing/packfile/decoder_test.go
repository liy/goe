package packfile

import (
	"testing"

	"github.com/liy/goe/tests"
)

func TestDecode(t *testing.T) {
	fixture := tests.GetFixture("topo-sort")
	idxFile := fixture.IndexFile("../../repos")
	packFile := fixture.PackFile("../../repos")
	
	pack, err := Decode(packFile, idxFile)
	if err != nil {
		t.Fatal(err)
	}

	tests.ToMatchSnapshot(t, pack)
}