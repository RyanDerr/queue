package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptions(t *testing.T) {
	t.Parallel()

	t.Run("Default Options", func(t *testing.T) {
		t.Parallel()
		opts := GetOpts()

		require.NotNil(t, opts)
		assert.Empty(t, opts.withMsg)
	})

	t.Run("WithMessage", func(t *testing.T) {
		t.Parallel()
		opts := GetOpts(WithMsg("custom error message"))

		require.NotNil(t, opts)
		assert.Equal(t, "custom error message", opts.withMsg)
	})
}
