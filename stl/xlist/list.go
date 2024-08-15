package xlist

type Node[E any] struct {
	next, prev *Node[E]

	list  *BLoop[E]
	Value E
}

func (n *Node[E]) Next() *Node[E] {
	if p := n.next; n.list != nil && p != &n.list.root {
		return p
	}
	return nil
}

func (n *Node[E]) Prev() *Node[E] {
	if p := n.prev; n.list != nil && p != &n.list.root {
		return p
	}
	return nil
}

// This is a doubly linked list ring
type BLoop[E any] struct {
	root Node[E]
	len  int
}

func New[E any]() *BLoop[E] {
	b := new(BLoop[E])
	b.root.next = &b.root
	b.root.prev = &b.root
	b.len = 0
	return b
}

func (b *BLoop[E]) Len() int {
	return b.len
}
