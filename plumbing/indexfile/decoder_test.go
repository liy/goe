package indexfile

import (
	"testing"

	"github.com/liy/goe/tests"
)

func TestDecoder(t *testing.T) {
	fixture := tests.GetFixture("topo-sort")
	file := fixture.IndexFile("../../repos")
	idx, err := Decode(file)
	if err != nil {
		t.Fatal(err)
	}

	tests.ToMatchSnapshot(t, idx)
}
