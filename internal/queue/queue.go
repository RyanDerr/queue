package queue

import (
	"sync"

	"github.com/RyanDerr/queue/internal/errors"
	"github.com/RyanDerr/queue/internal/node"
	"github.com/RyanDerr/queue/internal/util"
)

// queue is a struct that represents a queue data structure. It contains
// a mutex for thread safety, the size of the queue, and pointers to the
// head and tail nodes of the queue.
type queue[T any] struct {
	sync.RWMutex

	// size keeps track of the number of elements currently in the queue.
	size uint
	// capacity is an optional limit on the number of elements the queue can hold.
	// If capacity is 0, the queue has no size limit.
	capacity uint
	// head is a pointer to the dummy head node of the queue.
	//  The actual first element in the queue is the next node after the head.
	head *node.Node[T]
	// tail is a pointer to the dummy tail node of the queue.
	// The actual last element in the queue is the previous node before the tail.
	tail *node.Node[T]
}

// New creates and returns a new instance of a queue.
func New[T any](opts ...Option) *queue[T] {
	const op = "queue.New"
	o := getOpts(opts...)
	head, tail := node.New(*new(T)), node.New(*new(T))

	if err := head.SetNext(tail); err != nil {
		panic(errors.Wrap(op, err))
	}

	if err := tail.SetPrev(head); err != nil {
		panic(errors.Wrap(op, err))
	}

	return &queue[T]{
		size:     0,
		head:     head,
		tail:     tail,
		capacity: o.capacity,
	}
}

// Len returns the number of elements currently in the queue.
func (q *queue[T]) Len() uint {
	q.RLock()
	defer q.RUnlock()
	return q.size
}

// Enqueue adds an element to the back of the queue. If the provided value is
// nil, it returns an error.
func (q *queue[T]) Enqueue(v T) error {
	const op = "queue.(Queue).Enqueue"
	q.Lock()
	defer q.Unlock()

	// If a capacity limit is set on the queue
	// and the queue has reached that limit, return an error.
	if q.capacity > 0 && q.size >= q.capacity {
		return errors.Wrap(op, ErrQueueFull)
	}

	newNode := node.New(v)

	// Get the node that is currently at the end of the queue, which is the node
	// that will become the new node's previous node in the queue.
	curLast, err := q.tail.GetPrev()
	switch {
	case err != nil:
		return errors.Wrap(op, err)
	case util.IsNil(curLast):
		return errors.Wrap(op, ErrInternal, errors.WithMsg("fetching the tail node's previous node resulted in nil"))
	}

	// Set the new node to be the next node after the current last node in the queue.
	err = curLast.SetNext(newNode)
	if err != nil {
		return errors.Wrap(op, err)
	}

	// Link the new node back to the current last node and forward to the tail.
	err = newNode.SetPrev(curLast)
	if err != nil {
		return errors.Wrap(op, err)
	}

	err = newNode.SetNext(q.tail)
	if err != nil {
		return errors.Wrap(op, err)
	}

	// Set the new node to be the previous node before the tail node in the queue.
	err = q.tail.SetPrev(newNode)
	if err != nil {
		return errors.Wrap(op, err)
	}

	q.size++
	return nil
}

// Dequeue removes and returns the element at the front of the queue. If the
// queue is empty, it returns an error.
func (q *queue[T]) Dequeue() (T, error) {
	const op = "queue.(Queue).Dequeue"
	q.Lock()
	defer q.Unlock()
	if q.size == 0 {
		return *new(T), errors.Wrap(op, ErrEmptyQueue)
	}

	// Get the node that is currently at the head of the queue, which is the
	// node that will be dequeued.
	curHead, err := q.head.GetNext()
	switch {
	case err != nil:
		return *new(T), errors.Wrap(op, err)
	case util.IsNil(curHead):
		return *new(T), errors.Wrap(op, ErrInternal, errors.WithMsg("fetching the head node's next node resulted in nil"))
	}

	// Get the next node after the current head node, which will become the new head of the queue.
	newHead, err := curHead.GetNext()
	switch {
	case err != nil:
		return *new(T), errors.Wrap(op, err)
	case util.IsNil(newHead):
		return *new(T), errors.Wrap(op, ErrInternal, errors.WithMsg("fetching the next node of the head node resulted in nil"))
	}

	// Set the next node after the current head node to be the new head of the queue.
	err = q.head.SetNext(newHead)
	if err != nil {
		return *new(T), errors.Wrap(op, err)
	}

	err = newHead.SetPrev(q.head)
	if err != nil {
		return *new(T), errors.Wrap(op, err)
	}

	q.size--
	return curHead.GetValue()
}

// Peak returns the element at the front of the queue without removing it. If the
// queue is empty, it returns an error.
func (q *queue[T]) Peak() (T, error) {
	const op = "queue.(Queue).Peak"
	q.RLock()
	defer q.RUnlock()
	if q.size == 0 {
		return *new(T), errors.Wrap(op, ErrEmptyQueue)
	}

	curHead, err := q.head.GetNext()
	switch {
	case err != nil:
		return *new(T), errors.Wrap(op, err)
	case util.IsNil(curHead):
		return *new(T), errors.Wrap(op, ErrInternal, errors.WithMsg("fetching the head node's next node resulted in nil"))
	}

	return curHead.GetValue()
}

// Clear removes all elements from the queue, resetting it to an empty state.
func (q *queue[T]) Clear() {
	q.Lock()
	defer q.Unlock()

	// Reset the head and tail nodes to their initial state, with the head's next node
	// pointing to the tail and the tail's previous node pointing to the head.
	if err := q.head.SetNext(q.tail); err != nil {
		panic(errors.Wrap("queue.(Queue).Clear", err))
	}
	if err := q.tail.SetPrev(q.head); err != nil {
		panic(errors.Wrap("queue.(Queue).Clear", err))
	}
	q.size = 0
}

// DequeueBack removes and returns the element at the back of the queue. If the
// queue is empty, it returns an error.
func (q *queue[T]) DequeueBack() (T, error) {
	const op = "queue.(Queue).DequeueBack"
	q.Lock()
	defer q.Unlock()

	if q.size == 0 {
		return *new(T), errors.Wrap(op, ErrEmptyQueue)
	}

	// Get the node that is currently at the tail of the queue, which is the
	// node that will be dequeued.
	curTail, err := q.tail.GetPrev()
	switch {
	case err != nil:
		return *new(T), errors.Wrap(op, err)
	case util.IsNil(curTail):
		return *new(T), errors.Wrap(op, ErrInternal, errors.WithMsg("fetching the tail node's previous node resulted in nil"))
	}

	// Get the previous node before the current tail node, which will become the new tail of the queue.
	newTail, err := curTail.GetPrev()
	switch {
	case err != nil:
		return *new(T), errors.Wrap(op, err)
	case util.IsNil(newTail):
		return *new(T), errors.Wrap(op, ErrInternal, errors.WithMsg("fetching the previous node of the tail node resulted in nil"))
	}

	// Set the previous node before the current tail node to be the new tail of the queue.
	err = q.tail.SetPrev(newTail)
	if err != nil {
		return *new(T), errors.Wrap(op, err)
	}

	err = newTail.SetNext(q.tail)
	if err != nil {
		return *new(T), errors.Wrap(op, err)
	}

	q.size--
	return curTail.GetValue()
}
