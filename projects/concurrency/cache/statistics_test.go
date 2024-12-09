package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHitRate(t *testing.T) {
	t.Run("no value has been read", func(t *testing.T) {
		c := NewCache[string, any](3)

		c.Put("name", "John")
		c.Put("age", 30)
		got := c.GetStatistics().hitRate
		require.Equal(t, float64(0), got)

	})

	t.Run("all keys present in cache", func(t *testing.T) {
		c := NewCache[string, any](3)

		input := map[string]any{
			"name":  "Joe",
			"age":   30,
			"hobby": "swimming",
		}

		for key, value := range input {
			c.Put(key, value)
		}
		for key := range input {
			c.Get(key)
		}

		got := c.GetStatistics().hitRate

		require.Equal(t, 100.00, got)
	})

	t.Run("some keys not in cache", func(t *testing.T) {
		c := NewCache[string, any](3)

		input := map[string]any{
			"name":  "Joe",
			"age":   30,
			"hobby": "swimming",
		}

		for key, value := range input {
			c.Put(key, value)
		}
		c.Get("not in cache")
		c.Get("role")
		c.Get("name")

		got := c.GetStatistics().hitRate

		require.Equal(t, 33.33, got)
	})

}

func TestTotalReadsAndWrites(t *testing.T) {
	t.Run("no item has been deleted", func(t *testing.T) {
		c := NewCache[string, any](3)

		input := map[string]any{
			"name":  "Joe",
			"age":   30,
			"hobby": "swimming",
		}

		for key, value := range input {
			c.Put(key, value)
		}
		for key := range input {
			c.Get(key)
		}

		got := c.GetStatistics().totalReadsAndWrites

		require.Equal(t, 6, got)

	})

	t.Run("item evicted", func(t *testing.T) {
		c := NewCache[string, any](3)

		input := map[string]any{
			"name":  "Joe",
			"age":   30,
			"hobby": "swimming",
			"role":  "QA tester",
		}

		for key, value := range input {
			c.Put(key, value)
		}
		c.Get("name")
		got := c.GetStatistics().totalReadsAndWrites

		require.Equal(t, 5, got)

	})

	t.Run("item refreshed", func(t *testing.T) {
		c := NewCache[string, any](3)

		input := map[string]any{
			"name":  "Joe",
			"age":   30,
			"hobby": "swimming",
			"role":  "QA tester",
		}

		for key, value := range input {
			c.Put(key, value)
		}
		// refreshing name
		c.Put("name", "Anna")

		got := c.GetStatistics().totalReadsAndWrites

		require.Equal(t, 5, got)

	})

	t.Run("key not in cache", func(t *testing.T) {
		c := NewCache[string, any](3)

		input := map[string]any{
			"name":  "Joe",
			"age":   30,
			"hobby": "swimming",
			"role":  "QA tester",
		}

		for key, value := range input {
			c.Put(key, value)
		}
		// refreshing name
		c.Get("not in cache")

		got := c.GetStatistics().totalReadsAndWrites

		require.Equal(t, 5, got)

	})
}

func TestAverageReads(t *testing.T) {
	t.Run("average reads", func(t *testing.T) {
		c := NewCache[string, any](3)
		input := map[string]string{
			"name":  "Joe",
			"hobby": "traveling"}

		for key, value := range input {
			c.Put(key, value)
			c.Get(key)
		}
		require.Equal(t, 1.0, c.GetStatistics().averageReads)
	})

	t.Run("things currently in the cache", func(t *testing.T) {
		c := NewCache[string, any](3)
		input := map[string]string{
			"name":  "Joe",
			"hobby": "traveling",
			"role" : "Tester"}

		for key, value := range input {
			c.Put(key, value)
			c.Get(key)
		}
		c.Put("not-in-cache", "not-in-cache")
		require.Equal(t, 0.67, c.GetStatistics().averageReads)
	})

}
