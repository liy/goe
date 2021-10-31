package packfile

import (
	"encoding/binary"
	"io"
	"testing"

	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/store"
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

func BenchmarkReadObject(b *testing.B) {
	dotgit := store.NewDotGit("/topo-sort/.git", tests.Embeded{})
	pack, _ := dotgit.Pack("6179faab20f2d649a12fd52aab3c8d6e32b27dcd")
	packIndex, _ := dotgit.PackIndex("6179faab20f2d649a12fd52aab3c8d6e32b27dcd")
	packReader := NewPackReader(pack, packIndex)

	hashes := []plumbing.Hash{
		plumbing.ToHash("b39f0caec88b946e189d37a0a803e3d14d1399a4"),
		plumbing.ToHash("8aca3c14f5d8f2c40bf743a5bd68af6902e0026a"),
		plumbing.ToHash("ee8feeb79f6fcc59efe9c844d7c61b6cf4380b1f"),
		plumbing.ToHash("4854f1ce5137086767e71b4d3010db28bcd09c49"),
		plumbing.ToHash("5fd75f24dad1b1a785f41efebe4f075e5b0239f9"),
		plumbing.ToHash("6cc7b92f2f0b38e239ecb0745e25d774804365e1"),
		plumbing.ToHash("60bf463ccea705a5f22fb93fc61d4472d0d54cc3"),
		plumbing.ToHash("a790b0d579d164252a7b05d3743edbdf13847334"),
		plumbing.ToHash("f46f4be6e6132ebfadb70cd1eacefddba7b7110c"),
		plumbing.ToHash("405616a32d26161a8536ffaeca21b73dc9ca7fa3"),
		plumbing.ToHash("0ad552cb81f6955f223423c30857cb97b7706a7f"),
		plumbing.ToHash("7230da8e6c9b5173b27021f4df27fdf613eac14b"),
		plumbing.ToHash("8aca3c14f5d8f2c40bf743a5bd68af6902e0026a"),
		plumbing.ToHash("14200aa62633da89bdaaff5f0877585e93da5d78"),
	}

	for n := 0; n < b.N; n++ {
		packReader.ReadObject(hashes[n%len(hashes)])
	}
}