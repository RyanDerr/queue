package queue

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		options  []Option
		wantSize uint
		wantCap  uint
	}{
		{
			name:     "default options",
			options:  nil,
			wantSize: 0,
			wantCap:  0,
		},
		{
			name:     "custom capacity",
			options:  []Option{WithCapacity(10)},
			wantSize: 0,
			wantCap:  10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := New[int](tt.options...)
			require.NotNil(t, q)
			assert.Equal(t, tt.wantSize, q.size)
			assert.Equal(t, tt.wantCap, q.capacity)
		})
	}
}

func TestEnqueue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		enqueue  []int
		wantSize uint
	}{
		{
			name:     "single element",
			enqueue:  []int{1},
			wantSize: 1,
		},
		{
			name:     "multiple elements",
			enqueue:  []int{1, 2, 3, 4, 5},
			wantSize: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := New[int]()
			for _, v := range tt.enqueue {
				require.NoError(t, q.Enqueue(v))
			}
			assert.Equal(t, tt.wantSize, q.size)
		})
	}

	t.Run("returns ErrQueueFull when at capacity", func(t *testing.T) {
		t.Parallel()

		q := New[int](WithCapacity(2))
		require.NoError(t, q.Enqueue(1))
		require.NoError(t, q.Enqueue(2))

		err := q.Enqueue(3)
		require.ErrorIs(t, err, ErrQueueFull)
		assert.Equal(t, uint(2), q.size)
	})
}

func TestEnqueueLinkedStructure(t *testing.T) {
	t.Parallel()

	q := New[int]()
	require.NoError(t, q.Enqueue(10))
	require.NoError(t, q.Enqueue(20))

	// head -> node(10) -> node(20) -> tail
	first, _ := q.head.GetNext()
	v1, _ := first.GetValue()
	assert.Equal(t, 10, v1)

	second, _ := first.GetNext()
	v2, _ := second.GetValue()
	assert.Equal(t, 20, v2)

	// Backward link: node(20) -> node(10)
	back, _ := second.GetPrev()
	bv, _ := back.GetValue()
	assert.Equal(t, 10, bv)
}

func TestDequeue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		enqueue   []int
		dequeues  int
		wantOrder []int
		wantErr   error
		wantSize  uint
	}{
		{
			name:    "empty queue returns ErrEmptyQueue",
			enqueue: nil,
			wantErr: ErrEmptyQueue,
		},
		{
			name:      "FIFO order",
			enqueue:   []int{10, 20, 30},
			dequeues:  3,
			wantOrder: []int{10, 20, 30},
			wantSize:  0,
		},
		{
			name:      "partial dequeue",
			enqueue:   []int{10, 20, 30},
			dequeues:  1,
			wantOrder: []int{10},
			wantSize:  2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := New[int]()
			for _, v := range tt.enqueue {
				require.NoError(t, q.Enqueue(v))
			}

			if tt.wantErr != nil {
				_, err := q.Dequeue()
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			for i := range tt.dequeues {
				got, err := q.Dequeue()
				require.NoError(t, err)
				assert.Equal(t, tt.wantOrder[i], got)
			}
			assert.Equal(t, tt.wantSize, q.size)
		})
	}
}

func TestDequeueRestoresHeadTailLink(t *testing.T) {
	t.Parallel()

	q := New[int]()
	require.NoError(t, q.Enqueue(1))
	require.NoError(t, q.Enqueue(2))
	_, err := q.Dequeue()
	require.NoError(t, err)
	_, err = q.Dequeue()
	require.NoError(t, err)

	next, err := q.head.GetNext()
	require.NoError(t, err)
	assert.Equal(t, q.tail, next)
}

func TestDequeueLeftRestoresHeadTailLink(t *testing.T) {
	t.Parallel()

	q := New[int]()
	require.NoError(t, q.Enqueue(1))
	require.NoError(t, q.Enqueue(2))
	_, err := q.DequeueLeft()
	require.NoError(t, err)
	_, err = q.DequeueLeft()
	require.NoError(t, err)

	next, err := q.head.GetNext()
	require.NoError(t, err)
	assert.Equal(t, q.tail, next)
}

func TestPeak(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		enqueue   []int
		wantValue int
		wantErr   error
		wantSize  uint
	}{
		{
			name:    "empty queue returns ErrEmptyQueue",
			enqueue: nil,
			wantErr: ErrEmptyQueue,
		},
		{
			name:      "returns front without removing",
			enqueue:   []int{10, 20},
			wantValue: 10,
			wantSize:  2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := New[int]()
			for _, v := range tt.enqueue {
				require.NoError(t, q.Enqueue(v))
			}

			if tt.wantErr != nil {
				_, err := q.Peak()
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			got, err := q.Peak()
			require.NoError(t, err)
			assert.Equal(t, tt.wantValue, got)
			assert.Equal(t, tt.wantSize, q.size)

			got2, err := q.Peak()
			require.NoError(t, err)
			assert.Equal(t, got, got2)
		})
	}
}

func TestClear(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		enqueue []int
	}{
		{
			name:    "on empty queue",
			enqueue: nil,
		},
		{
			name:    "on populated queue",
			enqueue: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := New[int]()
			for _, v := range tt.enqueue {
				require.NoError(t, q.Enqueue(v))
			}
			q.Clear()

			assert.Equal(t, uint(0), q.size)

			next, err := q.head.GetNext()
			require.NoError(t, err)
			assert.Equal(t, q.tail, next)

			prev, err := q.tail.GetPrev()
			require.NoError(t, err)
			assert.Equal(t, q.head, prev)
		})
	}
}

func TestDequeueLeft(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		enqueue   []int
		dequeues  int
		wantOrder []int
		wantErr   error
		wantSize  uint
	}{
		{
			name:    "empty queue returns ErrEmptyQueue",
			enqueue: nil,
			wantErr: ErrEmptyQueue,
		},
		{
			name:      "LIFO order",
			enqueue:   []int{10, 20, 30},
			dequeues:  3,
			wantOrder: []int{30, 20, 10},
			wantSize:  0,
		},
		{
			name:      "partial dequeue",
			enqueue:   []int{10, 20, 30},
			dequeues:  1,
			wantOrder: []int{30},
			wantSize:  2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := New[int]()
			for _, v := range tt.enqueue {
				require.NoError(t, q.Enqueue(v))
			}

			if tt.wantErr != nil {
				_, err := q.DequeueLeft()
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			for i := range tt.dequeues {
				got, err := q.DequeueLeft()
				require.NoError(t, err)
				assert.Equal(t, tt.wantOrder[i], got)
			}
			assert.Equal(t, tt.wantSize, q.size)
		})
	}
}

func TestClearThenReuse(t *testing.T) {
	t.Parallel()

	q := New[int]()
	require.NoError(t, q.Enqueue(1))
	q.Clear()
	require.NoError(t, q.Enqueue(99))

	v, err := q.Dequeue()
	require.NoError(t, err)
	assert.Equal(t, 99, v)
}

func TestLen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		enqueue    int
		dequeue    int
		wantLength uint
	}{
		{
			name:       "empty",
			enqueue:    0,
			dequeue:    0,
			wantLength: 0,
		},
		{
			name:       "after enqueue",
			enqueue:    3,
			dequeue:    0,
			wantLength: 3,
		},
		{
			name:       "after partial dequeue",
			enqueue:    3,
			dequeue:    2,
			wantLength: 1,
		},
		{
			name:       "after draining",
			enqueue:    3,
			dequeue:    3,
			wantLength: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := New[int]()
			for i := range tt.enqueue {
				require.NoError(t, q.Enqueue(i))
			}
			for range tt.dequeue {
				_, err := q.Dequeue()
				require.NoError(t, err)
			}
			assert.Equal(t, tt.wantLength, q.Len())
		})
	}
}

func TestInterleavedOperations(t *testing.T) {
	t.Parallel()

	q := New[int]()
	require.NoError(t, q.Enqueue(1))
	require.NoError(t, q.Enqueue(2))

	v, _ := q.Dequeue()
	assert.Equal(t, 1, v)

	require.NoError(t, q.Enqueue(3))
	v, _ = q.Dequeue()
	assert.Equal(t, 2, v)

	v, _ = q.Dequeue()
	assert.Equal(t, 3, v)
}

func TestConcurrentAccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		prefill      int
		enqueue      int
		dequeue      int
		peak         int
		wantFinalLen uint
	}{
		{
			name:         "concurrent enqueue only",
			prefill:      0,
			enqueue:      100,
			dequeue:      0,
			wantFinalLen: 100,
		},
		{
			name:         "concurrent dequeue only",
			prefill:      50,
			enqueue:      0,
			dequeue:      50,
			wantFinalLen: 0,
		},
		{
			name:         "concurrent enqueue and dequeue",
			prefill:      50,
			enqueue:      50,
			dequeue:      50,
			wantFinalLen: 50,
		},
		{
			name:         "concurrent peak does not mutate",
			prefill:      1,
			enqueue:      0,
			dequeue:      0,
			peak:         50,
			wantFinalLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := New[int]()

			for i := range tt.prefill {
				require.NoError(t, q.Enqueue(i))
			}

			var wg sync.WaitGroup
			wg.Add(tt.enqueue + tt.dequeue + tt.peak)
			for i := range tt.enqueue {
				go func(v int) {
					defer wg.Done()
					_ = q.Enqueue(v)
				}(i)
			}
			for range tt.dequeue {
				go func() {
					defer wg.Done()
					_, _ = q.Dequeue()
				}()
			}
			for range tt.peak {
				go func() {
					defer wg.Done()
					_, _ = q.Peak()
				}()
			}
			wg.Wait()

			assert.Equal(t, tt.wantFinalLen, q.Len())
		})
	}
}
