package indexfile

import (
	"encoding/hex"
	"testing"

	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/store"
	"github.com/liy/goe/tests"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	idx *Index
}

func (suite *Suite) SetupTest() {
	suite.idx = NewIndex(tests.NewEmbeded(suite.T()).PackIndex())
}

func (suite *Suite) TestGetOffset() {
	hash := "0a2d236a6eb889bc433e712bab3e3a9b8889e988"
	offset, ok := suite.idx.GetOffset(plumbing.ToHash(hash))
	if !ok {
		suite.T().Fatalf("Cannot find hash %s", hash)
	}

	suite.Equalf(1481, int(offset), "Should give correct offset for hash %s", hash)
}

func (suite *Suite) TestGetHashFromOffset() {
	offset := 1481
	hash, ok := suite.idx.GetHashFromOffset(uint64(offset))
	if !ok {
		suite.T().Fatalf("Cannot find hash of offset %v", offset)
	}

	suite.Equal("0a2d236a6eb889bc433e712bab3e3a9b8889e988", hex.EncodeToString(hash), "Should give correct hash for offset %v", offset)
}

func (suite *Suite) TestGetHash() {
	expectedHash := "0a2d236a6eb889bc433e712bab3e3a9b8889e988"
	position := 1
	data := suite.idx.GetHash(position)
	hash := hex.EncodeToString(data)

	suite.Equal(expectedHash, hash, "second(position %v) object is %s", position, expectedHash)
}

func TestSuite(t *testing.T) {
    suite.Run(t, new(Suite))
}

var idx *Index
func init() {
	dotgit := store.NewDotGit("/topo-sort/.git", tests.Embeded{})
	file, _:= dotgit.PackIndex("6179faab20f2d649a12fd52aab3c8d6e32b27dcd")
	idx = NewIndex(file)
}

func BenchmarkGetOffset(b *testing.B) {
	hash := plumbing.ToHash("0a2d236a6eb889bc433e712bab3e3a9b8889e988")
	for n := 0; n < b.N; n++ {
		idx.GetOffset(hash)
	}
}

func BenchmarkGetHashFromOffset(b *testing.B) {
	offset := 1481
	for n := 0; n < b.N; n++ {
		idx.GetHashFromOffset(uint64(offset))
	}
}

func BenchmarkGetHash(b *testing.B) {
	position := 1
	for n := 0; n < b.N; n++ {
		idx.GetHash(position)
	}
}