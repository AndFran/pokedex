package pokecache

import (
	"sync"
	"time"
)

type PokeCache interface {
	Add(key string, value []byte)
	Get(key string) ([]byte, bool)
	Delete(key string)
	ViewKeys() []string
}

type Cache struct {
	cacheEntries map[string]CacheEntry
	Mu           *sync.Mutex
	Duration     time.Duration
}

type CacheEntry struct {
	CreatedAt time.Time
	Value     []byte
}

func NewMemoryCache(interval time.Duration) PokeCache {
	c := Cache{
		cacheEntries: make(map[string]CacheEntry),
		Mu:           &sync.Mutex{},
		Duration:     interval,
	}
	go c.reapLoop()
	return c
}

func (c Cache) reapLoop() {
	ticker := time.NewTicker(c.Duration)
	for _ = range ticker.C {
		for key, val := range c.cacheEntries {
			now := time.Now().UTC()
			expires := val.CreatedAt.Add(c.Duration)
			if now.After(expires) {
				c.Delete(key)
			}
		}
	}
}

func (c Cache) Add(key string, value []byte) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.cacheEntries[key] = CacheEntry{
		CreatedAt: time.Now().UTC(),
		Value:     value,
	}
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	val, ok := c.cacheEntries[key]
	return val.Value, ok
}

func (c Cache) Delete(key string) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	delete(c.cacheEntries, key)
}

func (c Cache) ViewKeys() []string {
	//not thread safe, items could be removed after the read etc.
	results := make([]string, 0, len(c.cacheEntries))
	for v, _ := range c.cacheEntries {
		results = append(results, v)
	}
	return results
}
