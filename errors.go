package queue

import (
	"github.com/RyanDerr/queue/internal/node"
	"github.com/RyanDerr/queue/internal/queue"
)

var (
	// ErrEmptyQueue is returned when an operation is attempted on an empty queue.
	ErrEmptyQueue = queue.ErrEmptyQueue
	// ErrNilNode is returned when an operation is performed on or with a nil node.
	ErrNilNode = node.ErrNilNode
	// ErrInitialization is returned when the queue fails to initialize properly.
	ErrInitialization = queue.ErrInitialization
	// ErrInternal is returned when unexpected internal behavior is encountered.
	ErrInternal = queue.ErrInternal
)
