package queue

// Option configures a queue at construction time.
type Option func(*options)

type options struct {
	capacity uint
}

func getOpts(opts ...Option) options {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// WithCapacity sets the maximum number of elements the queue can hold.
// Enqueue will return ErrQueueFull once this limit is reached.
// A value of 0 means unlimited.
func WithCapacity(size uint) Option {
	return func(o *options) {
		o.capacity = size
	}
}
