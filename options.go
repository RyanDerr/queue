package queue

import "github.com/RyanDerr/queue/internal/queue"

// Option configures a Queue at construction time.
type Option = queue.Option

// WithCapacity sets the maximum number of elements the queue can hold.
// Enqueue will return ErrQueueFull once this limit is reached.
// A value of 0 means unlimited.
var WithCapacity = queue.WithCapacity
