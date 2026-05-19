package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOpts(t *testing.T) {
	t.Parallel()

	t.Run("defaults to zero capacity", func(t *testing.T) {
		t.Parallel()
		o := getOpts()
		assert.Equal(t, uint(0), o.capacity)
	})

	t.Run("applies WithCapacity", func(t *testing.T) {
		t.Parallel()
		o := getOpts(WithCapacity(10))
		assert.Equal(t, uint(10), o.capacity)
	})
}
