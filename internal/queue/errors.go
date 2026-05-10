package queue

import "errors"

var (
	// ErrInitialization is an error that is returned when the queue fails to initialize properly.
	ErrInitialization = errors.New("queue failed to initialize")
	// ErrEmptyQueue is an error that is returned when an operation is attempted on an empty queue.
	ErrEmptyQueue = errors.New("queue is empty")
	// ErrInternal is a generic error that is returned when an unexpected behavior is encountered.
	ErrInternal = errors.New("encountered unexpected internal error")
)
