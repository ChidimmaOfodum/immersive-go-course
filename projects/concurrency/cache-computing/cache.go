package cachecomputing

import (
	"concurrency/cache"
	"sync"
)

type Creator[K comparable, V any] func(K) V

type ComputingCache[K comparable, V any] struct {
	cache                   cache.Cache[K, V]
	creator                 Creator[K, V]
	pendingKeys             map[K]chan V
	numberOfWaitingRoutines map[K]int
	mu                      sync.Mutex
}

func NewComputingCache[K comparable, V any](entryLimit int, creator Creator[K, V]) *ComputingCache[K, V] {
	return &ComputingCache[K, V]{
		cache:                   cache.NewCache[K, V](entryLimit),
		creator:                 creator,
		numberOfWaitingRoutines: map[K]int{},
		pendingKeys:             map[K]chan V{},
	}
}

func (c *ComputingCache[K, V]) Get(key K) V {
	value, present := c.cache.Get(key)
	if present {
		return *value
	}

	c.mu.Lock()
	if channel, ok := c.pendingKeys[key]; ok {
		c.numberOfWaitingRoutines[key]++
		c.mu.Unlock()
		created := <-channel
		return created
	}
	c.pendingKeys[key] = make(chan V)
	c.mu.Unlock()

	computedValue := c.creator(key)
	c.cache.Put(key, computedValue)

	// send value to all waiting routines
	c.mu.Lock()
	if count, ok := c.numberOfWaitingRoutines[key]; ok {
		channel := c.pendingKeys[key]
		c.mu.Unlock()
		for i := 0; i < count; i++ {
			channel <- computedValue
		}
		c.mu.Lock()
		c.numberOfWaitingRoutines[key] = 0
	}
	c.mu.Unlock()

	//delete from pending keys
	c.mu.Lock()
	delete(c.pendingKeys, key)
	c.mu.Unlock()

	return computedValue
}
