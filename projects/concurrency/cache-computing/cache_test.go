package cachecomputing

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// creator function
var callCount = 0

func creator(word string) int {
	callCount++
	return len(word)
}
func TestComputingCache(t *testing.T) {

	t.Run("test get", func(t *testing.T) {
		cache := NewComputingCache[string, int](10, creator)

		got := cache.Get("Thailand")
		require.Equal(t, 8, got)
	})

	t.Run("test concurrency", func(t *testing.T) {
		callCount = 0
		var wg sync.WaitGroup

		cache := NewComputingCache[string, int](10, creator)
		wg.Add(2)
		go func() {
			defer wg.Done()
			cache.Get("Singapore")
		}()
		go func() {
			defer wg.Done()
			cache.Get("Singapore")
		}()
		wg.Wait()
		require.Equal(t, 2, callCount)

	})
}
