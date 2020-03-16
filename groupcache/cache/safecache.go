package cache

import (
	"cache/lru"
	"sync")

// Impalement concurrent Cache
type cache struct {
	lru *lru.Cache
	mu sync.Mutex
	cacheBytes64 int64
}

func (c* cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Lazy Initialization
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes64,nil)
	}
	c.lru.Add(key,value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if value, ok := c.lru.Get(key); ok {
		return value.(ByteView), true
	}
	return
}

