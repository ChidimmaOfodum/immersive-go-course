package cachecomputing

import (
	"concurrency/cache"
	"sync"
)

type Creator[K comparable, V any] func(K) V
type Pending_key[V any] struct {
	pLock                   *sync.Mutex
	numberOfWaitingRoutines int
	channel                 chan V
}

type ComputingCache[K comparable, V any] struct {
	gLock sync.Mutex // this lock needs to held when checking if a key is present in the cache
	cache cache.Cache[K, V]

	creator Creator[K, V]

	mu          sync.Mutex // this lock needs to be held when reading/writing to pendingKeys map
	pendingKeys map[K]*Pending_key[V]
}

func NewComputingCache[K comparable, V any](entryLimit int, creator Creator[K, V]) *ComputingCache[K, V] {
	return &ComputingCache[K, V]{
		cache:       cache.NewCache[K, V](entryLimit),
		creator:     creator,
		pendingKeys: map[K]*Pending_key[V]{},
	}
}

func (c *ComputingCache[K, V]) Get(key K) V {
	value, present := c.cache.Get(key)
	if present {
		return *value
	}

	c.gLock.Lock()
	// recheck whether key is present
	if value, present = c.cache.Get(key); present {
		c.gLock.Unlock()
		return *value
	}
	c.gLock.Unlock()

	c.mu.Lock()
	if value, ok := c.pendingKeys[key]; ok {
		value.numberOfWaitingRoutines++
		c.mu.Unlock()
		created := <-value.channel
		return created
	}
	var pendingKeyLock sync.Mutex
	c.pendingKeys[key] = &Pending_key[V]{&pendingKeyLock, 0, make(chan V)}
	c.mu.Unlock()

	computedValue := c.creator(key)
	c.cache.Put(key, computedValue)

	// send value to all waiting routines

	c.mu.Lock()
	pendingKey := c.pendingKeys[key]
	delete(c.pendingKeys, key)
	c.mu.Unlock()

	pendingKey.pLock.Lock()
	if pendingKey.numberOfWaitingRoutines > 0 {
		channel := pendingKey.channel
		for i := 0; i < pendingKey.numberOfWaitingRoutines; i++ {
			channel <- computedValue
		}
	}
	pendingKey.pLock.Unlock()
	return computedValue
}
