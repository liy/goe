package goe

import (
	"testing"

	"github.com/liy/goe/plumbing"
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

func TestGetCommit(t *testing.T) {
	c, err := repo.GetCommit(plumbing.ToHash("4854f1ce5137086767e71b4d3010db28bcd09c49"))
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, len(c.Parents) == 1, "only has 1 parent")
	assert.Equal(t, "c53c4c18e245d880899405c07eb4d01b735b72ad", c.Parents[0].String(), "has correct parent commit")
	assert.Equal(t, "Fix package keywords.\n", c.Message, "commit has correct message")
}

func TestGetAnnotatedTag(t *testing.T) {
}
