package concurrency

import (
	"math"
	"time"
)

type Statistics struct {
	hitRate             float64
	unReadEntries      int
	averageHit          float64
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
	unReadEntries   int
	totalReads      int
	totalWrites     int
	successfulReads int
}

func NewCache[K comparable, V any](entryLimit int) Cache[K, V] {
	return Cache[K, V]{entryLimit: entryLimit, cache: make(map[K]cacheEntry[V])}
}

func (c *Cache[K, V]) Put(key K, value V) bool {
	c.totalWrites += 1
	entry := cacheEntry[V]{value: value, timeStamp: time.Now()}

	if _, ok := c.cache[key]; ok {
		c.cache[key] = entry
		return ok
	}

	//inserting new item
	if len(c.cache) == c.entryLimit {
		lastUsed := getLastRecentlyUsedEntry(c.cache)
		if (c.cache[lastUsed].numberOfTimesRead == 0) {
			c.unReadEntries += 1
		}
		delete(c.cache, lastUsed)
	}
	c.cache[key] = entry
	return false
}

func (c *Cache[K, V]) Get(key K) (*V, bool) {
	c.totalReads += 1
	entry, ok := c.cache[key]
	if !ok {
		return nil, ok
	}
	entry.timeStamp = time.Now()
	entry.numberOfTimesRead++
	c.cache[key] = entry // ToDO 
	c.successfulReads += 1
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
		averageHit:          roundFloat(averageReads, 2),
		totalReadsAndWrites: c.totalReads + c.totalWrites,
		unReadEntries:      c.unReadEntries + getUnReadEntries(c.cache),
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
