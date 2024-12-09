package cache

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
			"name":       "John",
			"age":        30,
			"hobby":      "swimming",
			"friend":     "Amy",
			"occupation": "Pilot",
		}
		for key, value := range input {
			c.Put(key, value)
		}
		require.Equal(t, 3, len(c.cache))
		c.Get("name")
	})
}

func TestLastRecentlyUsedEntry(t *testing.T) {
	c := make(map[string]cacheEntry[string, string])
	c["1"] = cacheEntry[string, string]{value: "entry1", timeStamp: time.Now()}
	c["2"] = cacheEntry[string, string]{value: "entry2", timeStamp: time.Now().Add(10 * time.Second)}
	c["3"] = cacheEntry[string, string]{value: "entry3", timeStamp: time.Now().Add(20 * time.Second)}
	c["4"] = cacheEntry[string, string]{value: "entry4", timeStamp: time.Now().Add(-30 * time.Second)}

	key := getLastRecentlyUsedEntry(c)
	require.Equal(t, "4", key)
}

func TestGetUnReadEntries(t *testing.T) {
	tests := map[string]struct {
		input    map[string]cacheEntry[string, string]
		expected int
	}{
		"all entries read": {
			input: map[string]cacheEntry[string, string]{
				"name":     {value: "Anna", timeStamp: time.Now(), numberOfTimesRead: 2},
				"hobby":    {value: "swimming", timeStamp: time.Now(), numberOfTimesRead: 1},
				"lastName": {value: "Joe", timeStamp: time.Now(), numberOfTimesRead: 3},
			},
			expected: 0},
		"no entries read": {
			input: map[string]cacheEntry[string, string]{
				"name":     {value: "Anna", timeStamp: time.Now(), numberOfTimesRead: 0},
				"hobby":    {value: "swimming", timeStamp: time.Now(), numberOfTimesRead: 0},
				"lastName": {value: "Joe", timeStamp: time.Now(), numberOfTimesRead: 0},
			},
			expected: 3},

		"some entries read": {
			input: map[string]cacheEntry[string, string]{
				"name":     {value: "Anna", timeStamp: time.Now(), numberOfTimesRead: 2},
				"hobby":    {value: "swimming", timeStamp: time.Now(), numberOfTimesRead: 0},
				"lastName": {value: "Joe", timeStamp: time.Now(), numberOfTimesRead: 1},
			},
			expected: 1},
	}

	for key, value := range tests {
		t.Run(key, func(t *testing.T) {
			got := unReadEntriesInMap[string, string](value.input)
			require.Equal(t, value.expected, got)
		})
	}
}

