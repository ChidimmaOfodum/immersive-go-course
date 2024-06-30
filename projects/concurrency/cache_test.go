package concurrency

import (
	"testing"
	"time"
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
		require.Equal(t, "John", val.value)
		require.True(t, ok)
	})

	t.Run("get item from a cache", func(t *testing.T) {
		c := NewCache[string, any](3)

		val, ok := c.Get("name")
		require.False(t, ok)
		require.Nil(t, val)

		c.Put("name", "John")

		val, ok = c.Get("name")
		require.True(t, ok)
		require.NotNil(t, val)
		require.Equal(t, "John", *val)
	})

	t.Run("does not exceed entryLimit", func(t *testing.T) {
		c := NewCache[string, any](3)
		input := map[string]any{
			"name":   "John",
			"age":    30,
			"hobby":  "swimming",
			"friend": "Amy",
		}
		for key, value := range input {
			c.Put(key, cacheEntry[any]{value: value, readTime: time.Now()})
		}
		require.Equal(t, 3, len(c.cache))
	})

	t.Run("test hit rate", func (t *testing.T)  {
		c := NewCache[string, any](3)

		for i:=0; i < 20; i++ {
			c.Get("name")
		}
		require.Equal(t, c.statistics.hitRate, 0)
		c.Put("name", "John")

		for i:=0; i < 10; i++ {
			c.Get("name")
		}
		require.Equal(t, 10, c.statistics.hitRate)
	})
}


func TestLastRecentlyUsedEntry(t *testing.T) {
	c := make(map[string]cacheEntry[string])
	c["secondEntry"] = cacheEntry[string]{value: "Hello", readTime: time.Now()}
	c["firstEntry"] = cacheEntry[string]{value: "Hello2", readTime: time.Now().Add(10 * time.Second)}
	c["thirdEntry"] = cacheEntry[string]{value: "Hello3", readTime: time.Now().Add(20 * time.Second)}

	key := getLastRecentlyUsedEntry(c)
	require.Equal(t, "secondEntry", key)
}
