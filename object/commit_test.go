package object

import (
	"testing"

	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/tests"
)

func TestDecodeCommit(t *testing.T) {
	hash := "f2010ee942a47bec0ca7e8f04240968ea5200735"
	commit, err := DecodeCommit(plumbing.GetRawObject(t, hash))
	if err != nil {
		t.Fatal(err)
	}
	tests.ToMatchSnapshot(t, commit)
}