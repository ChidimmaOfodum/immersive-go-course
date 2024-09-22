package concurrency

import (
	"math"
	"time"
	"sync"
)

type Statistics struct {
	hitRate             float64
	unReadEntries      int
	averageReads          float64
	totalReadsAndWrites int
}

type cacheEntry[V any] struct {
	value     V
	timeStamp time.Time
	numberOfTimesRead int
}

type Cache[K comparable, V any] struct {
	entryLimit      int
	cache           map[K]cacheEntry[V]
	unReadEvictedEntries   int
	totalReads      int
	totalWrites     int
	successfulReads int
	mu sync.Mutex
}

func NewCache[K comparable, V any](entryLimit int) Cache[K, V] {
	return Cache[K, V]{entryLimit: entryLimit, cache: make(map[K]cacheEntry[V])}
}

func (c *Cache[K, V]) Put(key K, value V) bool {
	c.mu.Lock()
	c.totalWrites += 1
	entry := cacheEntry[V]{value: value, timeStamp: time.Now()}

	if value, ok := c.cache[key]; ok {
		entry.numberOfTimesRead = value.numberOfTimesRead
		c.cache[key] = entry
		c.mu.Unlock()
		return ok
	}

	//inserting new item
	if len(c.cache) == c.entryLimit {
		lastUsed := getLastRecentlyUsedEntry(c.cache)
		if (c.cache[lastUsed].numberOfTimesRead == 0) {
			c.unReadEvictedEntries += 1
		}
		delete(c.cache, lastUsed)
	}
	c.cache[key] = entry
	c.mu.Unlock()
	return false
}

func (c *Cache[K, V]) Get(key K) (*V, bool) {
	c.mu.Lock()
	c.totalReads += 1
	entry, ok := c.cache[key]
	if !ok {
		c.mu.Unlock()
		return nil, ok
	}
	entry.timeStamp = time.Now()
	entry.numberOfTimesRead++
	c.cache[key] = entry // ToDO 
	c.successfulReads += 1
	c.mu.Unlock()
	return &entry.value, ok
}

func (c *Cache[K, V]) GetStatistics() Statistics {
	var hitRate float64 
	if c.totalReads != 0 {
		result := float64(c.successfulReads)/float64(c.totalReads) * 100
		hitRate = roundFloat(result, 2)
	}
	//calculate average reads
	totalReadsInCache := 0
	for _, value := range c.cache {
		totalReadsInCache += value.numberOfTimesRead
	}
	averageReads := float64(totalReadsInCache) / float64(len(c.cache))

	return Statistics{
		hitRate: hitRate,
		averageReads:          roundFloat(averageReads, 2),
		totalReadsAndWrites: c.totalReads + c.totalWrites,
		unReadEntries:      c.unReadEvictedEntries + unReadEntriesInMap(c.cache),
	}
}

func getLastRecentlyUsedEntry[K comparable, V any](entries map[K]cacheEntry[V]) K {
	var lastUsedKey K
	var readTime = time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC)

	for key, value := range entries {
			if value.timeStamp.Before(readTime) {
				readTime = value.timeStamp
				lastUsedKey = key
			}
		}
	return lastUsedKey
}

func unReadEntriesInMap[K comparable, V any](entries map[K]cacheEntry[V]) int {
	count := 0

	for _, value := range entries {
		if value.numberOfTimesRead == 0 {
			count++
		}
	}
	return count
}

func roundFloat(num float64, precision int) float64 {
	multiplier := math.Pow(10.0, float64(precision))
	return math.Round(num * multiplier) / multiplier
}
