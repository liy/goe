package packfile

import (
	"testing"

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