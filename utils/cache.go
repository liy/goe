package utils

import (
	"container/list"
	"sync"

	"github.com/liy/goe/plumbing"
)

type CacheItem interface {
	Hash() plumbing.Hash
	Size() int64
}

type Cache interface {
	Add(newItem CacheItem) bool
	Get(key plumbing.Hash) (CacheItem, bool)
	Contains(key plumbing.Hash) bool
}

type LRU struct {
	// Max allowed memroy footprint
	MaxSize int64
	// Current memmory footprint
	Size int64
	lst *list.List
	store map[plumbing.Hash]*list.Element

	mutex sync.Mutex
}

func NewLRU(maxSize int64) Cache{
	return &LRU{
		MaxSize: maxSize,
		lst: list.New(),
		store: make(map[plumbing.Hash]*list.Element, 2000),
	}
}

func (cache *LRU) Add(newItem CacheItem) bool {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// Do not allow item larger than cache size 
	if newItem.Size() > cache.MaxSize {
		return false
	}

	// update a existing cache item
	if le, ok := cache.store[newItem.Hash()]; ok {
		cache.lst.MoveToFront(le)
		cache.Size += newItem.Size() - le.Value.(CacheItem).Size()
	} else { // add new cache item
		le := cache.lst.PushFront(newItem)
		cache.store[newItem.Hash()] = le
		cache.Size += newItem.Size()
	}

	for cache.Size > cache.MaxSize {
		el := cache.lst.Back()
		if el == nil {
			cache.Size = 0
			break
		}

		objectToRemove := el.Value.(CacheItem)
		cache.Size -= objectToRemove.Size()
		cache.lst.Remove(el)
		delete(cache.store, objectToRemove.Hash())
	}
	return true
}

func (cache *LRU) Get(hash plumbing.Hash) (CacheItem, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	el, ok := cache.store[hash]
	if !ok {
		return nil, false
	}
	cache.lst.MoveToFront(el)
	return el.Value.(CacheItem), true
}

func (cache *LRU) Contains(hash plumbing.Hash) bool {
	_, ok := cache.store[hash]
	return ok
}