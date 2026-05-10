package queue

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	q := New[int]()
	require.NotNil(t, q)
	assert.Equal(t, uint(0), q.Len())
}

func TestEnqueue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		enqueue    []int
		wantLength uint
	}{
		{
			name:       "single element",
			enqueue:    []int{1},
			wantLength: 1,
		},
		{
			name:       "multiple elements",
			enqueue:    []int{1, 2, 3, 4, 5},
			wantLength: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := New[int]()
			for _, v := range tt.enqueue {
				require.NoError(t, q.Enqueue(v))
			}
			assert.Equal(t, tt.wantLength, q.Len())
		})
	}
}

func TestDequeue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		enqueue   []int
		dequeues  int
		wantOrder []int
		wantErr   error
	}{
		{
			name:    "empty queue returns ErrEmptyQueue",
			enqueue: nil,
			wantErr: ErrEmptyQueue,
		},
		{
			name:      "FIFO order",
			enqueue:   []int{1, 2, 3},
			dequeues:  3,
			wantOrder: []int{1, 2, 3},
		},
		{
			name:      "partial dequeue",
			enqueue:   []int{10, 20, 30},
			dequeues:  1,
			wantOrder: []int{10},
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
		})
	}
}

func TestPeak(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		enqueue    []int
		wantValue  int
		wantErr    error
		wantLength uint
	}{
		{
			name:    "empty queue returns ErrEmptyQueue",
			enqueue: nil,
			wantErr: ErrEmptyQueue,
		},
		{
			name:       "returns front without removing",
			enqueue:    []int{1, 2, 3},
			wantValue:  1,
			wantLength: 3,
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
			assert.Equal(t, tt.wantLength, q.Len())

			// Peak must be idempotent
			got2, _ := q.Peak()
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
		{name: "on empty queue", enqueue: nil},
		{name: "on populated queue", enqueue: []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := New[int]()
			for _, v := range tt.enqueue {
				require.NoError(t, q.Enqueue(v))
			}
			q.Clear()
			assert.Equal(t, uint(0), q.Len())
		})
	}
}

func TestClearThenReuse(t *testing.T) {
	t.Parallel()

	q := New[int]()
	require.NoError(t, q.Enqueue(1))
	require.NoError(t, q.Enqueue(2))
	q.Clear()

	require.NoError(t, q.Enqueue(99))
	got, err := q.Dequeue()
	require.NoError(t, err)
	assert.Equal(t, 99, got)
}

func TestLen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		enqueue    int
		dequeue    int
		wantLength uint
	}{
		{name: "empty", enqueue: 0, dequeue: 0, wantLength: 0},
		{name: "after enqueue", enqueue: 3, dequeue: 0, wantLength: 3},
		{name: "after partial dequeue", enqueue: 3, dequeue: 2, wantLength: 1},
		{name: "after draining", enqueue: 3, dequeue: 3, wantLength: 0},
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

func TestGenericTypes(t *testing.T) {
	t.Parallel()

	t.Run("string", func(t *testing.T) {
		t.Parallel()
		q := New[string]()
		require.NoError(t, q.Enqueue("hello"))
		got, err := q.Dequeue()
		require.NoError(t, err)
		assert.Equal(t, "hello", got)
	})

	t.Run("float64", func(t *testing.T) {
		t.Parallel()
		q := New[float64]()
		require.NoError(t, q.Enqueue(3.14))
		got, err := q.Dequeue()
		require.NoError(t, err)
		assert.Equal(t, 3.14, got)
	})

	t.Run("struct", func(t *testing.T) {
		t.Parallel()
		type item struct {
			id   int
			name string
		}
		q := New[item]()
		want := item{1, "test"}
		require.NoError(t, q.Enqueue(want))
		got, err := q.Dequeue()
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

// TestConcurrentAccess uses explicit goroutines because t.Parallel() only runs
// independent subtests concurrently — it cannot exercise concurrent access to a
// single shared queue instance, which is required to validate the mutex-based
// thread safety of the queue implementation under contention.
func TestConcurrentAccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		prefill      int
		enqueue      int
		dequeue      int
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
			name:         "concurrent enqueue and dequeue",
			prefill:      50,
			enqueue:      50,
			dequeue:      50,
			wantFinalLen: 50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := New[int]()

			for i := range tt.prefill {
				require.NoError(t, q.Enqueue(i))
			}

			// Goroutines are necessary here to create genuine concurrent contention
			// on the shared queue. The race detector (-race) validates correctness.
			var wg sync.WaitGroup
			wg.Add(tt.enqueue + tt.dequeue)
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
			wg.Wait()

			assert.Equal(t, tt.wantFinalLen, q.Len())
		})
	}
}

func TestQueueInterface(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   func() any
	}{
		{name: "int", fn: func() any { return New[int]() }},
		{name: "string", fn: func() any { return New[string]() }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.NotNil(t, tt.fn())
		})
	}
}
