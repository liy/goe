package indexfile

import (
	"testing"

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