package list

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInsertAtHeadOfList(t *testing.T) {
	t.Run("insert at head", func(t *testing.T) {
		list := &LinkedList[string]{}
		list.InsertAtHead(&Node[string]{Key: "name"})
		list.InsertAtHead(&Node[string]{Key: "hobby"})
		list.InsertAtHead(&Node[string]{Key: "address"})
		list.InsertAtHead(&Node[string]{Key: "role"})

		expect := "role\naddress\nhobby\nname\n"

		got := list.printList()

		require.Equal(t, expect, got)
		require.NotNil(t, list.Tail)
		require.Equal(t, "name", list.Tail.Key)
	})
}

func TestMoveNodeToHead(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		list_empty := &LinkedList[string]{}

		list_empty.MoveToHead(&Node[string]{Key: "role"})
		got := list_empty.printList()
		require.Equal(t, "role\n", got)
	})

	t.Run("target node at tail", func(t *testing.T) {
		list := &LinkedList[string]{}
		node3 := &Node[string]{Key: "node3"}

		list.InsertAtHead(node3)
		list.InsertAtHead(&Node[string]{Key: "node2"})
		list.InsertAtHead(&Node[string]{Key: "node1"})

		got := list.printList()
		require.Equal(t, "node1\nnode2\nnode3\n", got)

		list.MoveToHead(node3)
		got = list.printList()
		require.Equal(t, "node3\nnode1\nnode2\n", got)
		require.NotNil(t, list.Tail)
		require.Equal(t, "node2", list.Tail.Key)
	})

	t.Run("target node at middle", func(t *testing.T) {
		list := &LinkedList[string]{}
		node2 := &Node[string]{Key: "node2"}

		list.InsertAtHead(&Node[string]{Key: "node3"})
		list.InsertAtHead(node2)
		list.InsertAtHead(&Node[string]{Key: "node1"})

		got := list.printList()
		require.Equal(t, "node1\nnode2\nnode3\n", got)

		list.MoveToHead(node2)
		got = list.printList()
		require.Equal(t, "node2\nnode1\nnode3\n", got)
		require.NotNil(t, list.Tail)
		require.Equal(t, "node3", list.Tail.Key)
	})

	t.Run("target node at beginning", func(t *testing.T) {
		list := &LinkedList[string]{}
		node1 := &Node[string]{Key: "node1"}

		list.InsertAtHead(&Node[string]{Key: "node3"})
		list.InsertAtHead(&Node[string]{Key: "node2"})
		list.InsertAtHead(node1)

		got := list.printList()
		require.Equal(t, "node1\nnode2\nnode3\n", got)

		list.MoveToHead(node1)
		got = list.printList()
		require.Equal(t, "node1\nnode2\nnode3\n", got)
		require.NotNil(t, list.Tail)
		require.Equal(t, "node3", list.Tail.Key)
	})
}

func TestDeleteLastNode(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		list_empty := &LinkedList[string]{}
		list_empty.DeleteLastNode()
		got := list_empty.printList()

		require.Equal(t, "", got)
	})

	t.Run("one item in list", func(t *testing.T) {
		node := &Node[string]{Key: "node1"}
		list := &LinkedList[string]{Head: node, Tail: node}

		list.DeleteLastNode()
		got := list.printList()

		require.Equal(t, "", got)
	})

	t.Run("multiple items in list", func(t *testing.T) {
		list := &LinkedList[string]{}
		list.InsertAtHead(&Node[string]{Key: "node3"})
		list.InsertAtHead(&Node[string]{Key: "node2"})
		list.InsertAtHead(&Node[string]{Key: "node1"})

		list.DeleteLastNode()
		got := list.printList()

		require.Equal(t, "node1\nnode2\n", got)
		require.Equal(t, list.Tail.Key, "node2")
	})
}
