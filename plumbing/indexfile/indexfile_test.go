package indexfile

import (
	"encoding/hex"
	"testing"

	"github.com/liy/goe/plumbing"
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
