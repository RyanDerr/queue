package queue

import "github.com/RyanDerr/queue/internal/queue"

// Queue is an interface that defines the methods for a queue data structure. It
// includes methods for getting the length of the queue, adding an element to the
// back of the queue, removing and returning the element at the front of the
// queue, returning the element at the front of the queue without removing it, and
// clearing all elements from the queue.
type Queue[T any] interface {
	// Len returns the number of elements in the queue.
	Len() uint
	// Enqueue adds an element to the back of the queue.
	Enqueue(v T) error
	// Dequeue removes and returns the element at the front of the queue.
	Dequeue() (T, error)
	// DequeueLeft removes and returns the element at the back of the queue.
	DequeueLeft() (T, error)
	// Peak returns the element at the front of the queue without removing it.
	Peak() (T, error)
	// Clear removes all elements from the queue.
	Clear()
}

// New creates and returns a new instance of a Queue.
func New[T any](opts ...Option) Queue[T] {
	return queue.New[T](opts...)
}
