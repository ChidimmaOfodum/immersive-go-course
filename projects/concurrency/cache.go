package concurrency

type Cache[K comparable, V any] struct {
	entryLimit int
	cache      map[K]V
}

func NewCache[K comparable, V any](entryLimit int) Cache[K, V] {
	return Cache[K, V]{entryLimit: entryLimit, cache: make(map[K]V, entryLimit)}
}

func (c *Cache[K, V]) Put(key K, value V) bool {
	_, ok := c.cache[key]
	c.cache[key] = value
	return ok
}

func (c * Cache[K, V]) Get(key K) (*V, bool) {
	val, ok := c.cache[key]
	if !ok {
		return nil, ok
	}
	return &val, ok
}
