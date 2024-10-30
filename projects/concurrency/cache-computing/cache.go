package cachecomputing

import (
	"concurrency/cache"
	"fmt"
	"os"

	// "os"
	"sync"
	"time"
)

type Creator[K comparable, V any] func(K) V

type ComputingCache[K comparable, V any] struct {
	cache       cache.Cache[K, V]
	creator     Creator[K, V]
	pendingKeys map[K]struct{}
	mu          sync.Mutex
}

func NewComputingCache[K comparable, V any](entryLimit int, creator Creator[K, V]) *ComputingCache[K, V] {
	return &ComputingCache[K, V]{
		cache:       cache.NewCache[K, V](entryLimit),
		creator:     creator,
		pendingKeys: map[K]struct{}{},
	}

}

func (c *ComputingCache[K, V]) Get(key K) V {
	value, present := c.cache.Get(key)
	if present {
		return *value
	}

	c.mu.Lock()
	if _, ok := c.pendingKeys[key]; ok {
		for {
			if created, _ok := c.cache.Get(key); _ok {
				defer c.mu.Unlock()
				return *created
			}
			fmt.Fprintf(os.Stderr, "Waiting for %v to be calculated...\n", key)
			c.mu.Unlock()
			time.Sleep(time.Second * 1)
			c.mu.Lock()
		}
	}
	c.pendingKeys[key] = struct{}{}
	c.mu.Unlock()

	computedValue := c.creator(key)

	c.mu.Lock()
	delete(c.pendingKeys, key)
	c.mu.Unlock()

	c.cache.Put(key, computedValue)
	return computedValue
}
