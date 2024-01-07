package pokecache

import (
	"sync"
	"time"
)

type cache struct {
	cache  map[string]cacheEntry
	mux    sync.RWMutex
	ticker time.Ticker
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *cache {
	cache := cache{
		cache:  make(map[string]cacheEntry),
		mux:    sync.RWMutex{},
		ticker: *time.NewTicker(interval),
	}

	// call the reapLoop every `interval` seconds
	go func() {
		<-cache.ticker.C
		cache.reapLoop(interval)
	}()

	return &cache
}

func (c *cache) Add(key string, val []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *cache) Get(key string) (val []byte, exists bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	cacheEntry, exists := c.cache[key]
	return cacheEntry.val, exists
}

func (c *cache) reapLoop(interval time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()
	for url, ce := range c.cache {
		if time.Now().Sub(ce.createdAt) > interval {
			delete(c.cache, url)
		}
	}
}
