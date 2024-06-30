package concurrency

import (
	"time"
)

type Statistics struct {
	hitRate        int
	unReadEndtries int
	averageHit     float32
	readsAndWrite  int
}

type cacheEntry[V any]struct {
	value V
	readTime time.Time
}

type Cache[K comparable, V any] struct {
	entryLimit int
	cache      map[K]cacheEntry[V]
	totalHits int
	statistics Statistics
}

func NewCache[K comparable, V any](entryLimit int) Cache[K, V] {
	return Cache[K, V]{entryLimit: entryLimit, cache: make(map[K]cacheEntry[V])}
}

func (c *Cache[K, V]) Put(key K, value V) bool {
	_, ok := c.cache[key]
	// replacing value of an existing key
	if ok {
		c.cache[key] = cacheEntry[V]{value: value, readTime: time.Now()}
		return ok
	}

	//inserting new item
	if len(c.cache) == c.entryLimit {
		lastUsed := getLastRecentlyUsedEntry(c.cache)
		delete(c.cache, lastUsed)
	}
	c.cache[key] = cacheEntry[V]{value: value, readTime: time.Now()}
	return ok
}

func (c *Cache[K, V]) Get(key K) (*V, bool) {
	c.totalHits += 1
	entry, ok := c.cache[key]
	if !ok {
		return nil, ok
	}
	entry.readTime = time.Now()
	c.statistics.hitRate += 1
	return &entry.value, ok
}

func (c *Cache[K, V]) GetStatistics() Statistics {
	return c.statistics
}

func getLastRecentlyUsedEntry[K comparable, V any](entries map[K]cacheEntry[V]) K {
	var lastUsedKey K
	var readTime time.Time
	count := 1

	for key, value := range entries {
		if count == 1 {
			readTime = value.readTime
			lastUsedKey = key
			count++
			continue
		} else {
			if value.readTime.Before(readTime) {
				readTime = value.readTime
				lastUsedKey = key
			}
		}
	}
	return lastUsedKey

}
