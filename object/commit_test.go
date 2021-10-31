package object

import (
	"testing"

	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/tests"
)

func TestDecodeCommit(t *testing.T) {
	hash := "f2010ee942a47bec0ca7e8f04240968ea5200735"
	file := tests.NewEmbeded(t).GetObjectFile(hash)
	raw := plumbing.NewRawObject(plumbing.ToHash(hash))
	err := raw.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	commit, err := DecodeCommit(raw)
	if err != nil {
		t.Fatal(err)
	}
	tests.ToMatchSnapshot(t, commit)
}