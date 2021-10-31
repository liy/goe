package packfile

import (
	"encoding/binary"
	"io"
	"testing"

	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/tests"
	"github.com/stretchr/testify/suite"
)

type PackFileSuite struct {
	suite.Suite
	packReader *PackReader
}

func (suite *PackFileSuite) SetupTest() {
	fixture := tests.NewEmbeded(suite.T())
	suite.packReader = NewPackReader(fixture.Pack(), fixture.PackIndex())
}

func (suite *PackFileSuite) TestReadObject() {
	hash := "b39f0caec88b946e189d37a0a803e3d14d1399a4"
	raw, err := suite.packReader.ReadObject(plumbing.ToHash(hash))
	if err != nil {
		suite.T().Fatal(err)
	}

	expectedContent := "tree 7637310b6bd9226d5425a0f1585b22b69050eadb\nparent 4d507ff6044d78b6809f240333bae341f508f06b\nauthor liy <liy8272@gmail.com> 1413661783 +0100\ncommitter liy <liy8272@gmail.com> 1413661783 +0100\n\nUpdate ignore and package info\n"

	suite.Equal("b39f0caec88b946e189d37a0a803e3d14d1399a4", raw.Hash().String(), "object hash is correct")
	suite.Equal(plumbing.OBJ_COMMIT, raw.Type, "object is commit type")
	suite.Equal(225, int(raw.DeflatedSize), "object has correct defalated size")
	suite.Equal(expectedContent, string(raw.Data), "object has correct content")
}

func (suite *PackFileSuite) TestPackReader() {
	pack := new(Pack)
	// Skip signature
	suite.packReader.Seek(4, io.SeekStart)
	// Version
	data := make([]byte, 4)
	suite.packReader.Read(data)
	pack.Version = int32(binary.BigEndian.Uint32(data))

	suite.Equal(2, int(pack.Version), "Only support version 2")
}

func TestSuite(t *testing.T) {
    suite.Run(t, new(PackFileSuite))
}
