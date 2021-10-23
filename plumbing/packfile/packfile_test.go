package packfile

import (
	"encoding/binary"
	"io"
	"testing"

	"github.com/liy/goe/fixtures"
	"github.com/liy/goe/plumbing"
	"github.com/stretchr/testify/assert"
)

var packReader *PackReader

func init() {
	fixture := fixtures.NewRepositoryFixture("topo-sort")
	packReader = NewPackReader(fixture.GetPackFilePath("../../repos"))
}

func TestReadObject(t *testing.T) {
	assert := assert.New(t)
	hash := "b39f0caec88b946e189d37a0a803e3d14d1399a4"
	raw, err := packReader.ReadObject(plumbing.ToHash(hash))
	if err != nil {
		t.Fatal(err)
	}

	expectedContent := "tree 7637310b6bd9226d5425a0f1585b22b69050eadb\nparent 4d507ff6044d78b6809f240333bae341f508f06b\nauthor liy <liy8272@gmail.com> 1413661783 +0100\ncommitter liy <liy8272@gmail.com> 1413661783 +0100\n\nUpdate ignore and package info\n"

	assert.Equal("b39f0caec88b946e189d37a0a803e3d14d1399a4", raw.Hash().String(), "object hash is correct")
	assert.Equal(plumbing.OBJ_COMMIT, raw.Type, "object is commit type")
	assert.Equal(225, int(raw.DeflatedSize), "object has correct defalated size")
	assert.Equal(expectedContent, string(raw.Data), "object has correct content")
}

func TestPackReader(t *testing.T) {
	pack := new(Pack)
	assert := assert.New(t)
	// Skip signature
	packReader.Seek(4, io.SeekStart)
	// Version
	data := make([]byte, 4)
	packReader.Read(data)
	pack.Version = int32(binary.BigEndian.Uint32(data))

	assert.Equal(2, int(pack.Version), "Only support version 2")
}
