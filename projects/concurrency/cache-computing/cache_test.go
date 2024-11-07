package cachecomputing

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// creator function

func TestComputingCache(t *testing.T) {

	t.Run("computes the value for a key when not in cache", func(t *testing.T) {
		var wg sync.WaitGroup
		var got int

		creator := func(word string) int {
			time.Sleep(time.Millisecond * 10)
			return len(word)
		}

		cache := NewComputingCache[string, int](10, creator)
		wg.Add(1)
		go func() {
			defer wg.Done()
			got = cache.Get("Singapore")
		}()
		wg.Wait()
		require.Equal(t, 9, got)
	})
	t.Run("calls creator only once", func(t *testing.T) {
		var callCount atomic.Int64
		var wg sync.WaitGroup

		creator := func(word string) int {
			callCount.Add(1)
			time.Sleep(time.Millisecond * 10)
			return len(word)
		}

		cache := NewComputingCache[string, int](10, creator)
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
			defer wg.Done()
			cache.Get("Singapore")
		}()
		}
		wg.Wait()
		require.Equal(t, int64(1), callCount.Load())
	})

	t.Run("waiting key should receive valid signal when the value has been computed", func(t *testing.T) {
		var callCount atomic.Int64
		var wg sync.WaitGroup
		var value1, value2, value3, value4, value5, value6 int

		creator := func(word string) int {
			callCount.Add(1)
			time.Sleep(time.Second * 2)
			return len(word)
		}
		cache := NewComputingCache[string, int](10, creator)
		wg.Add(6)
		go func() {
			defer wg.Done()
			value1 = cache.Get("Singapore")
		}()
		go func() {
			defer wg.Done()
			value2 = cache.Get("Cuba")
		}()
		go func() {
			defer wg.Done()
			value3 = cache.Get("Cuba")
		}()

		go func() {
			defer wg.Done()
			value4 = cache.Get("Singapore")
		}()

		go func() {
			defer wg.Done()
			value5 = cache.Get("Singapore")
		}()

		go func() {
			defer wg.Done()
			value6 = cache.Get("Cuba")
		}()
		wg.Wait()

		require.Equal(t, 9, value1)
		require.Equal(t, 4, value2)
		require.Equal(t, 4, value3)
		require.Equal(t, 9, value4)
		require.Equal(t, 9, value5)
		require.Equal(t, 4, value6)
	})
}
