package queue

import "github.com/RyanDerr/queue/internal/queue"

var (
	// ErrEmptyQueue is returned when an operation is attempted on an empty queue.
	ErrEmptyQueue = queue.ErrEmptyQueue
	// ErrInternal is returned when unexpected internal behavior is encountered.
	ErrInternal = queue.ErrInternal
)
