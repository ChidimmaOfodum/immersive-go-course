package list
import "fmt"

type Node[K comparable] struct {
	Key  K
	prev *Node[K]
	next *Node[K]
}

type LinkedList[K comparable] struct {
	Head *Node[K]
	Tail *Node[K]
}

func (list *LinkedList[K]) MoveToHead(node *Node[K]) {
	if node == list.Head {
		return
	}
	if list.Head == nil {
		list.InsertAtHead(node)
		return
	}

	if node == list.Tail {
		node.prev.next = nil
		list.Tail = node.prev
	} else {
		node.prev.next = node.next
		node.next.prev = node.prev
	}

	node.prev = nil
	node.next = list.Head
	list.Head.prev = node
	list.Head = node
}

func (list *LinkedList[K]) InsertAtHead(node *Node[K]) {
	if list.Head == nil {
		list.Head = node
		list.Tail = node
		return
	}

	list.Head.prev = node
	node.next = list.Head
	list.Head = node
}

func (list *LinkedList[K]) DeleteLastNode() {
	if list.Tail == nil {
		return
	}
	if list.Head == list.Tail {
		list.Head = nil
		list.Tail = nil
		return
	}
	tailNode := list.Tail
	tailNode.prev.next = nil
	list.Tail = tailNode.prev
	tailNode.prev = nil
}

func (list *LinkedList[K]) printList() string {
	var result string
	current := list.Head
	for current != nil {
		result += fmt.Sprintf("%v\n", current.Key)
		current = current.next
	}
	return result
}
