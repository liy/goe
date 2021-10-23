package indexfile

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/liy/goe/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestDecoder(t *testing.T) {
	fixture := fixtures.NewRepositoryFixture("topo-sort")
	file := fixture.IndexFile("../../repos")
	idx, err := Decode(file)
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadFile("./indexfile_snapshot.json")
	if err != nil {
		t.Fatal(err)
	}

	snapshot := new(Index)
	err = json.Unmarshal(data, &snapshot)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, snapshot.Fanout, idx.Fanout, "decoder should give correct fanout table")
	assert.Equal(t, snapshot.Buckets, idx.Buckets, "decoder should give correct buckets")
	assert.Equal(t, snapshot.Hashes, idx.Hashes, "decoder should give correct hashes slice")
	assert.Equal(t, snapshot.CRC, idx.CRC, "decoder should give correct CRC slice")
	assert.Equal(t, snapshot.Offset32, idx.Offset32, "decoder should give correct Offset32")
	assert.Equal(t, snapshot.Offset64, idx.Offset64, "decoder should give correct Offset64")
	assert.Equal(t, int(snapshot.Version), int(idx.Version), "decoder should give correct version")
	assert.Equal(t, snapshot.ReverseHash, idx.ReverseHash, "decoder should give correct reverse hash")
	assert.Equal(t, int(snapshot.NumObjects), int(idx.NumObjects), "decoder should give correct number of object")
}
