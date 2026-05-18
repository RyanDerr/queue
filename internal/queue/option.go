package queue

type Option func(*Options)

type Options struct {
	// withQueueSize specifies the maximum number of elements the queue can hold.
	withQueueSize uint
}

func GetOpts(opt ...Option) Options {
	opts := getDefaultOptions()
	for _, o := range opt {
		o(&opts)
	}
	return opts
}

func getDefaultOptions() Options {
	return Options{}
}

// WithQueueSize sets the maximum size of the queue.
// If the queue reaches this size, it will not accept
// new elements until some are removed.
func WithQueueSize(size uint) Option {
	return func(o *Options) {
		o.withQueueSize = size
	}
}
