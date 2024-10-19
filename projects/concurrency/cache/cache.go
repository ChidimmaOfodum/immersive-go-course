package cache

import (
	"concurrency/list"
	"math"
	"sync"
	"time"
)
type Statistics struct {
	hitRate             float64
	unReadEntries       int
	averageReads        float64
	totalReadsAndWrites int
}

type cacheEntry[K comparable, V any] struct {
	value             V
	timeStamp         time.Time
	numberOfTimesRead int
	node              *list.Node[K]
}

type Cache[K comparable, V any] struct {
	entryLimit int

	mu                   sync.Mutex
	cache                map[K]cacheEntry[K, V]
	unReadEvictedEntries int
	totalReads           int
	totalWrites          int
	successfulReads      int
	list                 *list.LinkedList[K]
}

func NewCache[K comparable, V any](entryLimit int) Cache[K, V] {
	newList := list.LinkedList[K]{}
	return Cache[K, V]{entryLimit: entryLimit, cache: make(map[K]cacheEntry[K, V]), list: &newList}
}

func (c *Cache[K, V]) Put(key K, value V) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.totalWrites += 1

	entry := cacheEntry[K, V]{value: value, timeStamp: time.Now()}

	if value, ok := c.cache[key]; ok {
		entry.numberOfTimesRead = value.numberOfTimesRead
		entry.node = value.node
		c.list.MoveToHead(value.node)
		c.cache[key] = entry
		return ok
	}

	if len(c.cache) == c.entryLimit {
		lastUsed := c.list.Tail
		if c.cache[lastUsed.Key].numberOfTimesRead == 0 {
			c.unReadEvictedEntries += 1
		}
		delete(c.cache, lastUsed.Key)
		c.list.DeleteLastNode()
	}

	node := &list.Node[K]{Key: key}
	entry.node = node
	c.cache[key] = entry
	c.list.InsertAtHead(node)
	return false
}

func (c *Cache[K, V]) Get(key K) (*V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.totalReads += 1
	entry, ok := c.cache[key]
	if !ok {
		return nil, ok
	}
	entry.timeStamp = time.Now()
	entry.numberOfTimesRead++
	c.list.MoveToHead(entry.node)
	c.cache[key] = entry
	c.successfulReads += 1
	return &entry.value, ok
}

func (c *Cache[K, V]) GetStatistics() Statistics {
	c.mu.Lock()
	defer c.mu.Unlock()
	var hitRate float64
	if c.totalReads != 0 {
		result := float64(c.successfulReads) / float64(c.totalReads) * 100
		hitRate = roundFloat(result, 2)
	}
	//calculate average reads
	totalReadsInCache := 0
	for _, value := range c.cache {
		totalReadsInCache += value.numberOfTimesRead
	}
	averageReads := float64(totalReadsInCache) / float64(len(c.cache))

	return Statistics{
		hitRate:             hitRate,
		averageReads:        roundFloat(averageReads, 2),
		totalReadsAndWrites: c.totalReads + c.totalWrites,
		unReadEntries:       c.unReadEvictedEntries + unReadEntriesInMap(c.cache),
	}
}

func getLastRecentlyUsedEntry[K comparable, V any](entries map[K]cacheEntry[K, V]) K {
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

func unReadEntriesInMap[K comparable, V any](entries map[K]cacheEntry[K, V]) int {
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
	return math.Round(num*multiplier) / multiplier
}
