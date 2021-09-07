package indexfile

import (
	"bytes"
	"encoding/binary"

	"github.com/liy/goe/plumbing"
)

const MSB_MASK_32 = uint32(1) << 31

type Index struct {
	Fanout   [256]uint32
	Buckets  [256]uint32
	Hashes   []plumbing.Hash
	CRC      [][4]byte
	Offset32 [][4]byte
	Offset64 []byte
	Version uint32
}

func (idx *Index) GetOffset(position int) int64 {
	// Mask out the MSB to construct offset
	offset := binary.BigEndian.Uint32(idx.Offset32[position][:]) & ^MSB_MASK_32
	// If the MSB is set, the rest 31bit offset is actually the index to the Offset64 bytes
	if idx.Offset32[position][0]&128 > 0 {
		return int64(binary.BigEndian.Uint64(idx.Offset64[offset : offset+8]))
	} else {
		return int64(offset)
	}
}

// search for offset from object hash
/**
1. Extract msb byte which gives the index to the fanout table
2. Read the value from the fanout table using the index gives the high
3. Read the value from the buckets table using the same index gives the low: low = high - value
4. Binary search the Hashes table start with the low and high value
*/
func (idx *Index) GetPosition(hash plumbing.Hash) int {
	i := hash[0]

	if idx.Buckets[i] == 0 {
		return -1
	}

	high := idx.Fanout[i]
	low := high - idx.Buckets[i]
	for {
		mid := (high + low) >> 1
		if low >= high {
			return -1
		}

		h := idx.Hashes[mid]
		r := bytes.Compare(hash[:], h[:])
		if r == 0 {
			return int(mid)	
		} else if r == -1 {
			high = mid
		} else {
			low = mid + 1
		}
	}
	
}
