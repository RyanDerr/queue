package node

import (
	"github.com/RyanDerr/queue/internal/errors"
	"github.com/RyanDerr/queue/internal/util"
)

// Node is a struct that represents a single element within a queue.
// It contains the value of the element, as well as pointers to it's
// previous and next nodes in the queue.
type Node[T any] struct {
	value T
	prev  *Node[T]
	next  *Node[T]
}

// New creates and returns a new instance of a Node with the provided value.
func New[T any](v T) *Node[T] {
	return &Node[T]{value: v}
}

// GetValue returns the value stored in the node. If the node is
// nil, it returns an error.
func (n *Node[T]) GetValue() (T, error) {
	const op = "node.(Node).GetValue"
	if util.IsNil(n) {
		return *new(T), errors.Wrap(op, ErrNilNode)
	}
	return n.value, nil
}

// GetNext returns the next node in the queue. If the current node
// is nil or the next node is nil, it returns an error.
func (n *Node[T]) GetNext() (*Node[T], error) {
	const op = "node.(Node).GetNext"
	switch {
	case util.IsNil(n):
		return nil, errors.Wrap(op, ErrNilNode)
	case util.IsNil(n.next):
		return nil, errors.Wrap(op, ErrNilNode, errors.WithMsg("fetching next node resulted in nil"))
	default:
		return n.next, nil
	}
}

// GetPrev returns the previous node in the queue. If the current node
// is nil or the previous node is nil, it returns an error.
func (n *Node[T]) GetPrev() (*Node[T], error) {
	const op = "node.(Node).GetPrev"
	switch {
	case util.IsNil(n):
		return nil, errors.Wrap(op, ErrNilNode)
	case util.IsNil(n.prev):
		return nil, errors.Wrap(op, ErrNilNode, errors.WithMsg("fetching previous node resulted in nil"))
	default:
		return n.prev, nil
	}
}

// SetNext sets the next node in the queue. If the current node or the
// provided next node is nil, it returns an error.
func (n *Node[T]) SetNext(next *Node[T]) error {
	const op = "node.(Node).SetNext"
	switch {
	case util.IsNil(n):
		return errors.Wrap(op, ErrNilNode)
	case util.IsNil(next):
		return errors.Wrap(op, ErrNilNode, errors.WithMsg("provided next node is nil"))
	}

	n.next = next
	return nil
}

// SetPrev sets the previous node in the queue. If the current node or the
// provided previous node is nil, it returns an error.
func (n *Node[T]) SetPrev(prev *Node[T]) error {
	const op = "node.(Node).SetPrev"
	switch {
	case util.IsNil(n):
		return errors.Wrap(op, ErrNilNode)
	case util.IsNil(prev):
		return errors.Wrap(op, ErrNilNode, errors.WithMsg("provided previous node is nil"))
	}

	n.prev = prev
	return nil
}
