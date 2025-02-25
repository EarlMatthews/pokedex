package pokecache

import (
	"sync"
	"time"
)

	type cacheEntry struct {
		createdAt	time.Time
		val			[]byte
	}

	type Cache struct {
		entries map[string]cacheEntry
		mu	sync.Mutex
		interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	c :=  &Cache{
		entries: make(map[string]cacheEntry),
		interval : interval,
	}
	// Start background reaper
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, exists := c.entries[key]
	if !exists || time.Since(entry.createdAt) > c.interval {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(){
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.entries {
			if entry.createdAt.Add(c.interval).Before(now) {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock() 
	}

}