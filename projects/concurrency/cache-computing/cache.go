package cachecomputing

import (
	"concurrency/cache"
)

type Creator[K comparable, V any] func(K)V

type ComputingCache[K comparable, V any] struct {
	cache   cache.Cache[K, V]
	creator Creator[K, V]
}

func NewComputingCache[K comparable, V any](entryLimit int, creator Creator[K, V]) *ComputingCache[K, V] {
	return &ComputingCache[K, V]{
		cache:   cache.NewCache[K, V](entryLimit),
		creator: creator,
	}

}

func (c *ComputingCache[K, V]) Get(key K) V {
	value, present := c.cache.Get(key)
	if present {
		return *value
	}

	computedValue := c.creator(key)
	c.cache.Put(key, computedValue)
	return computedValue
}
