package utils

import (
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