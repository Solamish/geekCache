package lru

import (
	"container/list"
)

type Cache struct {
	maxBytes int64
	// nbytes is the used memory
	nbytes int64
	// ll works as a queue
	// front is queue's head
	// back is queue's tail
	ll    *list.List
	cache map[string]*list.Element
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// New is the Constructor of cache
func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get a key's value
// then, add value to the queue's tail
func (cache *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := cache.cache[key]; ok {
		cache.ll.MoveToBack(ele)
		kv := ele.Value.(*entry)
		return kv.value, ok
	}
	return
}

// RemoveOldest removes the oldest item
// follow the strategy of LRU(Least Recently Used)
func (cache *Cache) RemoveOldest() {
	ele := cache.ll.Front()
	if ele != nil {
		cache.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(cache.cache, kv.key)
		cache.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if cache.OnEvicted != nil {
			cache.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add adds a value to the cache
// if key exists, push it to the queue's tail
// if nbytes > maxBytes, delete the item with LRU
func (cache *Cache)Add(key string, value Value) {
	if ele, ok := cache.cache[key]; ok {
		cache.ll.MoveToBack(ele)
		kv := ele.Value.(*entry)
		kv.value = value
		return
	}
	ele := cache.ll.PushBack(&entry{key, value})
	cache.cache[key] = ele
	cache.nbytes += int64(len(key)) + int64(value.Len())

	for cache.maxBytes != 0 && cache.maxBytes < cache.nbytes {
		cache.RemoveOldest()
	}
}

// Len return the length of cache entries
func (cache *Cache)Len() int{
	return cache.ll.Len()
}
