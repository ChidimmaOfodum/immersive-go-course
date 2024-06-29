package concurrency

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {

	t.Run("initializes a cache", func(t *testing.T) {
		c := NewCache[string, any](3)
		require.Equal(t, c.entryLimit, 3)
		require.NotNil(t, c.cache)
	})

	t.Run("add item to a cache", func(t *testing.T) {
		c := NewCache[string, any](3)
		ok := c.Put("name", "John")
		require.False(t, ok)
		val, ok := c.cache["name"]
		require.Equal(t, val, "John")
		require.True(t, ok)
	})

	t.Run("get item from a cache", func(t *testing.T) {
		c := NewCache[string, any](3)

		val, ok := c.Get("name")
		require.False(t, ok)
		require.Nil(t, val)

		c.cache["name"] = "John"

		val, ok = c.Get("name")
		require.True(t, ok)
		require.NotNil(t, val)
		require.Equal(t, "John", *val)
	})
}
