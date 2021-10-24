package goe

import (
	"testing"

	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/utils"
	"github.com/stretchr/testify/assert"
)

var repo *Repository

func init() {
	var err error
	repo, err = OpenRepository("../repos/topo-sort")
	if err != nil {
		panic(err)
	}
}

func TestHead(t *testing.T) {
	ref, err := repo.HEAD()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "refs/heads/dev", ref.Target.ReferenceName(), "HEAD is pointed to dev branch")
}

func TestTryPeel(t *testing.T) {
	ref, err := repo.HEAD()
	if err != nil {
		t.Fatal(err)
	}

	targetRef := repo.TryPeel(ref.Target)
	assert.Equal(t, "f2010ee942a47bec0ca7e8f04240968ea5200735", targetRef.ReferenceName(), "HEAD pointed to dev branch tip commit")
}

func TestGetReferences(t *testing.T) {
	refs := repo.GetReferences()
	
	utils.ToMatchSnapshot(t, refs)
}

func TestGetCommit(t *testing.T) {
	c, err := repo.GetCommit(plumbing.ToHash("4854f1ce5137086767e71b4d3010db28bcd09c49"))
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, len(c.Parents) == 1, "only has 1 parent")
	assert.Equal(t, "c53c4c18e245d880899405c07eb4d01b735b72ad", c.Parents[0].String(), "has correct parent commit")
	assert.Equal(t, "Fix package keywords.\n", c.Message, "commit has correct message")
	utils.ToMatchSnapshot(t, c)
}

func TestGetCommits(t *testing.T) {
	cs, err := repo.GetCommits(plumbing.ToHash("c91773cc3da3c2c3c954626a0b6a44c3ac9e3e92"))
	if err != nil {
		t.Fatal(err)
	}
	
	utils.ToMatchSnapshot(t, cs)
}

func TestGetAnnotatedTag(t *testing.T) {
	tag, err := repo.GetAnnotatedTag(plumbing.ToHash("ca0a44b6eddd79547d1ad8bc94be987489edde2a"))
	if err != nil {
		t.Fatal(err)
	}

	utils.ToMatchSnapshot(t, tag)
}

func TestReadObject(t *testing.T) {
	obj, err := repo.ReadObject(plumbing.ToHash("4854f1ce5137086767e71b4d3010db28bcd09c49"))
	if err != nil {
		t.Fatal(err)
	}

	utils.ToMatchSnapshot(t, obj)
}