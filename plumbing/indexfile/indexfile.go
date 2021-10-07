package indexfile

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"os"

	"github.com/liy/goe/plumbing"
)

const MSB_MASK_32 = uint32(1) << 31

type Index struct {
	Fanout   [256]uint32
	Buckets  [256]uint32
	Hashes   []byte
	CRC      []byte
	Offset32 []byte
	Offset64 []byte
	Version uint32
	NumObjects uint32
	// Reverse maps offset to hash index
	ReverseHash map[uint64]uint32
}

func NewIndex(path string) *Index{
	// bytes, _ := ioutil.ReadFile(path)
	file, _ := os.Open(path)
	idx, _ := Decode(bufio.NewReader(file))
	return idx
}
 
func (idx *Index) getOffset(position int) int64 {
	// Mask out the MSB to construct offset
	i := position*4
	offset := binary.BigEndian.Uint32(idx.Offset32[i : i+4]) & ^MSB_MASK_32
	// If the MSB is set, the rest 31bit offset is actually the index to the Offset64 bytes
	if idx.Offset32[i]&128 > 0 {
		return int64(binary.BigEndian.Uint64(idx.Offset64[offset : offset+8]))
	} else {
		return int64(offset)
	}
}

func (idx *Index) GetHash(position int) []byte {
	i := position*20
	return idx.Hashes[i : i+20]
}

func (idx *Index) GetHashFromOffset(offset uint64) ([]byte, bool) {
	position, ok := idx.ReverseHash[offset]
	if !ok {
		return nil, false
	}
	return idx.GetHash(int(position)), true
}

// TODO: cache this
// search for offset from object hash
/**
1. Extract msb byte which gives the index to the fanout table
2. Read the value from the fanout table using the index gives the high
3. Read the value from the buckets table using the same index gives the low: low = high - value
4. Binary search the Hashes table start with the low and high value
*/
func (idx *Index) GetOffset(hash plumbing.Hash) (int64, bool) {
	i := hash[0]

	if idx.Buckets[i] == 0 {
		return 0, false
	}

	high := idx.Fanout[i]
	low := high - idx.Buckets[i]
	for {
		mid := (high + low) >> 1
		if low >= high {
			return 0, false
		}

		midIdx := mid*20
		r := bytes.Compare(hash[:], idx.Hashes[midIdx:midIdx+20])
		if r == 0 {
			return idx.getOffset(int(mid)), true	
		} else if r == -1 {
			high = mid
		} else {
			low = mid + 1
		}
	}
	
}
