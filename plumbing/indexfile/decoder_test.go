package indexfile

import (
	"testing"

	"github.com/liy/goe/fixtures"
	"github.com/liy/goe/utils"
)

func TestDecoder(t *testing.T) {
	fixture := fixtures.NewRepositoryFixture("topo-sort")
	file := fixture.IndexFile("../../repos")
	idx, err := Decode(file)
	if err != nil {
		t.Fatal(err)
	}

	utils.ToMatchSnapshot(t, idx)
}
