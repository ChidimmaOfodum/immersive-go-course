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
	statistics Statistics
}

func NewCache[K comparable, V any](entryLimit int) Cache[K, V] {
	return Cache[K, V]{entryLimit: entryLimit, cache: make(map[K]cacheEntry[V])}
}

func (c *Cache[K, V]) Put(key K, value V) bool {
	_, ok := c.cache[key]

	c.cache[key] = struct {
		value    V
		readTime time.Time
	}{value: value, readTime: time.Time{}}
	return ok
}

func (c *Cache[K, V]) Get(key K) (*V, bool) {
	entry, ok := c.cache[key]
	if !ok {
		return nil, ok
	}
	return &entry.value, ok
}

func (c *Cache[K, V]) GetStatistics() Statistics {
	return c.statistics
}
