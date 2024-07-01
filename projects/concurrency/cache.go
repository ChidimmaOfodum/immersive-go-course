package concurrency

import (
	"time"
)

type Statistics struct {
	hitRate        int
	unReadEndtries int
	averageHit     float32
	totalReadsAndWrites  int
}

type cacheEntry[V any]struct {
	value V
	timeStamp time.Time
	isRead bool
}

type Cache[K comparable, V any] struct {
	entryLimit int
	cache      map[K]cacheEntry[V]
	unReadEntries int
	averageHit float32
	totalReads int
	totalWrites int
	hitRate int
	totalHits int
}

func NewCache[K comparable, V any](entryLimit int) Cache[K, V] {
	return Cache[K, V]{entryLimit: entryLimit, cache: make(map[K]cacheEntry[V])}
}

func (c *Cache[K, V]) Put(key K, value V) bool {
	c.totalWrites += 1
	_, ok := c.cache[key]
	// replacing value of an existing key
	if ok {
		c.cache[key] = cacheEntry[V]{value: value, timeStamp: time.Now(), isRead: false}
		return ok
	}

	//inserting new item
	if len(c.cache) == c.entryLimit {
		lastUsed := getLastRecentlyUsedEntry(c.cache)
		if c.cache[lastUsed].isRead {
			c.unReadEntries += 1
		}
		delete(c.cache, lastUsed)
	}
	c.cache[key] = cacheEntry[V]{value: value, timeStamp: time.Now()}
	return ok
}

func (c *Cache[K, V]) Get(key K) (*V, bool) {
	c.totalReads += 1
	entry, ok := c.cache[key]
	if !ok {
		return nil, ok
	}
	entry.timeStamp = time.Now()
	c.hitRate += 1
	return &entry.value, ok
}

func (c *Cache[K, V]) GetStatistics() Statistics {
	return Statistics{
		hitRate: c.hitRate,
		averageHit: c.averageHit,
		totalReadsAndWrites: c.totalReads + c.totalWrites,
		unReadEndtries: c.unReadEntries + getUnReadEntries(c.cache),
	}
}

func getLastRecentlyUsedEntry[K comparable, V any](entries map[K]cacheEntry[V]) K {
	var lastUsedKey K
	var readTime time.Time
	count := 1

	for key, value := range entries {
		if count == 1 {
			readTime = value.timeStamp
			lastUsedKey = key
			count++
			continue
		} else {
			if value.timeStamp.Before(readTime) {
				readTime = value.timeStamp
				lastUsedKey = key
			}
		}
	}
	return lastUsedKey

}

func getUnReadEntries[K comparable, V any](entries map[K]cacheEntry[V]) int {
	count := 0

	for _, value := range entries {
		if !value.isRead {
			count++
		}
	}
	return count
}
