package utils

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"testing"

	"github.com/liy/goe/plumbing"
	"github.com/stretchr/testify/assert"
)

type ItemMock struct {
	hash plumbing.Hash
	size int64
}

func (it ItemMock) Hash() plumbing.Hash {
	return it.hash
}

func (it ItemMock) Size() int64 {
	return it.size
}

func TestAdd(t *testing.T) {
	itemMock := &ItemMock{
		plumbing.ToHash("a91773cc3da3c2c3c954626a0b6a44c3ac9e3e92"),
		// 20 bytes
		20,
	}
	
	// 10 bytes max size
	cache := NewLRU(10)
	success := cache.Add(itemMock)
	assert.True(t, len(cache.store) == 0, "should not be able add item whose size is larger than Cache.MaxSize")
	assert.False(t, success, "should add item should return fail")

	// update existing cache item
	cache = NewLRU(40)
	cache.Add(itemMock)
	cache.Add(itemMock)
	assert.True(t, cache.Size == 20, "cache current size should be updated")
	assert.True(t, len(cache.store) == 1, "should update exiting cache item")

	// Should remove the least used item (also first in first out) from the cache
	cache = NewLRU(40)
	cache.Add(&ItemMock{
		plumbing.ToHash("a91773cc3da3c2c3c954626a0b6a44c3ac9e3e92"),
		20,
	})
	cache.Add(&ItemMock{
		plumbing.ToHash("b91773cc3da3c2c3c954626a0b6a44c3ac9e3e92"),
		20,
	})
	cache.Add(&ItemMock{
		plumbing.ToHash("c91773cc3da3c2c3c954626a0b6a44c3ac9e3e92"),
		20,
	})
	assert.True(t, cache.Size == 40, "cache current size should be updated")
	assert.True(t, len(cache.store) == 2)
	assert.True(t, !cache.Contains(plumbing.ToHash("a91773cc3da3c2c3c954626a0b6a44c3ac9e3e92")) && cache.Contains(plumbing.ToHash("b91773cc3da3c2c3c954626a0b6a44c3ac9e3e92")) && cache.Contains(plumbing.ToHash("c91773cc3da3c2c3c954626a0b6a44c3ac9e3e92")), "least used item (first item) should be removed")
}

func TestGet(t *testing.T) {
	cache := NewLRU(40)
	cache.Add(&ItemMock{
		plumbing.ToHash("a91773cc3da3c2c3c954626a0b6a44c3ac9e3e92"),
		20,
	})
	item, ok := cache.Get(plumbing.ToHash("a91773cc3da3c2c3c954626a0b6a44c3ac9e3e92"))
	assert.True(t, ok, "should get the item successfully")
	assert.True(t, item.Hash().String() == "a91773cc3da3c2c3c954626a0b6a44c3ac9e3e92" , "should get the correct item")

	//  test Least recently used
	cache = NewLRU(40)
	cache.Add(&ItemMock{
		plumbing.ToHash("a91773cc3da3c2c3c954626a0b6a44c3ac9e3e92"),
		20,
	})
	cache.Add(&ItemMock{
		plumbing.ToHash("b91773cc3da3c2c3c954626a0b6a44c3ac9e3e92"),
		20,
	})
	cache.Get(plumbing.ToHash("a91773cc3da3c2c3c954626a0b6a44c3ac9e3e92"))

	cache.Add(&ItemMock{
		plumbing.ToHash("c91773cc3da3c2c3c954626a0b6a44c3ac9e3e92"),
		20,
	})
	ok = cache.Contains(plumbing.ToHash("b91773cc3da3c2c3c954626a0b6a44c3ac9e3e92"))
	assert.False(t, ok, "get least used item should fail")
}

func makeHashes(N int, seed int64) []plumbing.Hash {
	rand.Seed(seed)
	randomHash := func() plumbing.Hash {
		sha := sha1.New()
		sha.Write([]byte(fmt.Sprint(rand.Intn(N))))
		return plumbing.NewHash(sha.Sum(nil))
	}

	hashes := make([]plumbing.Hash, N)
	for i := 0; i < N; i++ {
		hashes[i] = randomHash()
	}

	return hashes
}

func shuffleHashes(hashes []plumbing.Hash) []plumbing.Hash {
	rand.Shuffle(len(hashes), func(i, j int) {
		hashes[i], hashes[j] = hashes[j], hashes[i]
	})

	return hashes
}

var numHashes = 100000
var addHashes  = makeHashes(numHashes, 2)
var getHashes  = shuffleHashes(makeHashes(numHashes, 2))
var defaultCache *LRU
func init() {
	defaultCache = NewLRU(40)
	for i := 0; i < numHashes; i++ {
		defaultCache.Add(&ItemMock{
			addHashes[i],
			5,
		})
	}
}

func BenchmarkAdd(b *testing.B) {
	cache := NewLRU(1024 * 4)
	for n := 0; n < b.N; n++ {
		cache.Add(&ItemMock{
			addHashes[n%numHashes],
			5,
		})
	}
}

func BenchmarkGet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		defaultCache.Get(getHashes[n%numHashes])
	}
}

func BenchmarkAddGet(b *testing.B) {
	cache := NewLRU(1024 * 4)
	for n := 0; n < b.N; n++ {
		cache.Add(&ItemMock{
			addHashes[n%numHashes],
			5,
		})
		
		if n%3 ==0 {
			cache.Get(getHashes[n%numHashes])
		}
	}
}