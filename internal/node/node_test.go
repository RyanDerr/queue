package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value int
	}{
		{
			name:  "positive",
			value: 42,
		},
		{
			name:  "zero",
			value: 0,
		},
		{
			name:  "negative",
			value: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			n := New(tt.value)
			require.NotNil(t, n)
			v, err := n.GetValue()
			require.NoError(t, err)
			assert.Equal(t, tt.value, v)
		})
	}
}

func TestNewNodeHasNilLinks(t *testing.T) {
	t.Parallel()

	n := New(1)
	_, err := n.GetNext()
	require.Error(t, err)
	_, err = n.GetPrev()
	require.Error(t, err)
}

func TestGetValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		node      *Node[int]
		wantValue int
		wantErr   error
	}{
		{
			name:      "returns stored value",
			node:      New(99),
			wantValue: 99,
		},
		{
			name:    "nil node returns ErrNilNode",
			node:    nil,
			wantErr: ErrNilNode,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v, err := tt.node.GetValue()
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantValue, v)
		})
	}
}

func TestSetNextAndGetNext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		receiver    *Node[int]
		arg         *Node[int]
		wantErr     error
		wantNextVal int
	}{
		{
			name:        "valid set and get",
			receiver:    New(1),
			arg:         New(2),
			wantNextVal: 2,
		},
		{
			name:     "nil receiver returns ErrNilNode",
			receiver: nil,
			arg:      New(1),
			wantErr:  ErrNilNode,
		},
		{
			name:     "nil argument returns ErrNilNode",
			receiver: New(1),
			arg:      nil,
			wantErr:  ErrNilNode,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.receiver.SetNext(tt.arg)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)

			got, err := tt.receiver.GetNext()
			require.NoError(t, err)
			v, _ := got.GetValue()
			assert.Equal(t, tt.wantNextVal, v)
		})
	}
}

func TestGetNextErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		node    *Node[int]
		wantErr error
	}{
		{name: "nil receiver", node: nil, wantErr: ErrNilNode},
		{name: "nil next field", node: New(1), wantErr: ErrNilNode},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := tt.node.GetNext()
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestSetPrevAndGetPrev(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		receiver    *Node[int]
		arg         *Node[int]
		wantErr     error
		wantPrevVal int
	}{
		{
			name:        "valid set and get",
			receiver:    New(2),
			arg:         New(1),
			wantPrevVal: 1,
		},
		{
			name:     "nil receiver returns ErrNilNode",
			receiver: nil,
			arg:      New(1),
			wantErr:  ErrNilNode,
		},
		{
			name:     "nil argument returns ErrNilNode",
			receiver: New(1),
			arg:      nil,
			wantErr:  ErrNilNode,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.receiver.SetPrev(tt.arg)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)

			got, err := tt.receiver.GetPrev()
			require.NoError(t, err)
			v, _ := got.GetValue()
			assert.Equal(t, tt.wantPrevVal, v)
		})
	}
}

func TestGetPrevErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		node    *Node[int]
		wantErr error
	}{
		{
			name:    "nil receiver",
			node:    nil,
			wantErr: ErrNilNode,
		},
		{
			name:    "nil prev field",
			node:    New(1),
			wantErr: ErrNilNode,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := tt.node.GetPrev()
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestChainTraversal(t *testing.T) {
	t.Parallel()

	a := New(1)
	b := New(2)
	c := New(3)

	require.NoError(t, a.SetNext(b))
	require.NoError(t, b.SetPrev(a))
	require.NoError(t, b.SetNext(c))
	require.NoError(t, c.SetPrev(b))

	// Forward: a -> b -> c
	n, err := a.GetNext()
	require.NoError(t, err)
	v, _ := n.GetValue()
	assert.Equal(t, 2, v)

	n, err = n.GetNext()
	require.NoError(t, err)
	v, _ = n.GetValue()
	assert.Equal(t, 3, v)

	// Backward: c -> b -> a
	n, err = c.GetPrev()
	require.NoError(t, err)
	n, err = n.GetPrev()
	require.NoError(t, err)
	v, _ = n.GetValue()
	assert.Equal(t, 1, v)
}
