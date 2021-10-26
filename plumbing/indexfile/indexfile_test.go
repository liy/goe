package indexfile

import (
	"encoding/hex"
	"testing"

	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/tests"
	"github.com/stretchr/testify/assert"
)

var idx *Index

func init() {
	fixture := tests.GetFixture("topo-sort")
	filePath := fixture.GetIndexFilePath("../../repos")
	idx = NewIndex(filePath)
}

func TestGetOffset(t *testing.T) {
	hash := "0a2d236a6eb889bc433e712bab3e3a9b8889e988"
	offset, ok := idx.GetOffset(plumbing.ToHash(hash))
	if !ok {
		t.Fatalf("Cannot find hash %s", hash)
	}

	assert.Equalf(t, 1481, int(offset), "Should give correct offset for hash %s", hash)
}

func TestGetHashFromOffset(t *testing.T) {
	offset := 1481
	hash, ok := idx.GetHashFromOffset(uint64(offset))
	if !ok {
		t.Fatalf("Cannot find hash of offset %v", offset)
	}

	assert.Equal(t, "0a2d236a6eb889bc433e712bab3e3a9b8889e988", hex.EncodeToString(hash), "Should give correct hash for offset %v", offset)
}

func TestGetHash(t *testing.T) {
	expectedHash := "0a2d236a6eb889bc433e712bab3e3a9b8889e988"
	position := 1
	data := idx.GetHash(position)
	hash := hex.EncodeToString(data)

	assert.Equal(t, expectedHash, hash, "second(position %v) object is %s", position, expectedHash)
}
