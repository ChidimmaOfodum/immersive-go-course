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
			"occupation": "Pilot",
		}
		for key, value := range input {
			c.Put(key, cacheEntry[string, any]{value: value, timeStamp: time.Now()})
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

// linked list operations

func TestInsertAtHeadOfList(t *testing.T) {
	t.Run("insert at head", func(t *testing.T) {
		list := &LinkedList[string]{}
		list.insertAtHead(&Node[string]{key: "name"})
		list.insertAtHead(&Node[string]{key: "hobby"})
		list.insertAtHead(&Node[string]{key: "address"})
		list.insertAtHead(&Node[string]{key: "role"})

		expect := "role\naddress\nhobby\nname\n"

		got := list.printList()

		require.Equal(t, expect, got)
		require.NotNil(t, list.tail)
		require.Equal(t, "name", list.tail.key)
	})
}

func TestMoveNodeToHead(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		list_empty := &LinkedList[string]{}

		list_empty.moveToHead(&Node[string]{key: "role"})
		got := list_empty.printList()
		require.Equal(t, "role\n", got)
	})

	t.Run("target node at tail", func(t *testing.T) {
		list := &LinkedList[string]{}
		node1 := &Node[string]{key: "node1"}
		node2 := &Node[string]{key: "node2", prev: node1}
		node3 := &Node[string]{key: "node3", prev: node2}
		list.tail = node3
		list.head = node1
		node1.next = node2
		node2.next = node3

		got := list.printList()
		require.Equal(t, "node1\nnode2\nnode3\n", got)

		list.moveToHead(node3)
		got = list.printList()
		require.Equal(t, "node3\nnode1\nnode2\n", got)
		require.NotNil(t, list.tail)
		require.Equal(t, "node2", list.tail.key)
	})

	t.Run("target node at middle", func(t *testing.T) {
		list := &LinkedList[string]{}
		node1 := &Node[string]{key: "node1"}
		node2 := &Node[string]{key: "node2", prev: node1}
		node3 := &Node[string]{key: "node3", prev: node2}
		list.tail = node3
		list.head = node1
		node1.next = node2
		node2.next = node3

		got := list.printList()
		require.Equal(t, "node1\nnode2\nnode3\n", got)

		list.moveToHead(node2)
		got = list.printList()
		require.Equal(t, "node2\nnode1\nnode3\n", got)
		require.NotNil(t, list.tail)
		require.Equal(t, "node3", list.tail.key)
	})

	t.Run("target node at beginning", func(t *testing.T) {
		list := &LinkedList[string]{}
		node1 := &Node[string]{key: "node1"}
		node2 := &Node[string]{key: "node2", prev: node1}
		node3 := &Node[string]{key: "node3", prev: node2}
		list.tail = node3
		list.head = node1
		node1.next = node2
		node2.next = node3

		got := list.printList()
		require.Equal(t, "node1\nnode2\nnode3\n", got)

		list.moveToHead(node1)
		got = list.printList()
		require.Equal(t, "node1\nnode2\nnode3\n", got)
		require.NotNil(t, list.tail)
		require.Equal(t, "node3", list.tail.key)
	})
}
