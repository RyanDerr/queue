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

func WithMsg(msg string) Option {
	return func(o *Options) {
		o.withMsg = msg
	}
}
