package cachecomputing

import (
	"sync/atomic"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// creator function

func TestComputingCache(t *testing.T) {
	t.Run("test concurrency", func(t *testing.T) {
		var callCount atomic.Int64
		var wg sync.WaitGroup

		creator := func(word string) int {
			callCount.Add(1)
			time.Sleep(time.Millisecond * 10)
			return len(word)
		}

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
		require.Equal(t, int64(1), callCount.Load())
	})
}
