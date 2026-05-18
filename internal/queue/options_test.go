package queue

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptions(t *testing.T) {
	t.Parallel()

	t.Run("Default Options", func(t *testing.T) {
		t.Parallel()
		opts := GetOpts()
		require.NotNil(t, opts)
		require.Zero(t, opts.withQueueSize)
	})

	t.Run("Custom Queue Size", func(t *testing.T) {
		t.Parallel()
		opts := GetOpts(WithQueueSize(10))
		require.NotNil(t, opts)
		require.Equal(t, uint(10), opts.withQueueSize)
	})
}
