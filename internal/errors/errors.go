package errors

import "fmt"

type Op string

func Wrap(op Op, e error, options ...Option) error {
	opts := GetOpts(options...)

	if opts.withMsg != "" {
		return fmt.Errorf("%s: %w: %s", op, e, opts.withMsg)
	}

	return fmt.Errorf("%s: %w", op, e)
}
