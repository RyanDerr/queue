package errors

type Option func(*Options)

type Options struct {
	withMsg string
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

// WithMsg allows you to provide a custom message that will be included in the error
// when using the Wrap function.
func WithMsg(msg string) Option {
	return func(o *Options) {
		o.withMsg = msg
	}
}
