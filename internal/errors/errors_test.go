package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrap(t *testing.T) {
	t.Parallel()

	baseErr := errors.New("base error")
	otherErr := errors.New("other error")

	tests := []struct {
		name       string
		op         Op
		err        error
		opts       []Option
		wantString string
		wantIs     error
	}{
		{
			name:       "wraps with op only",
			op:         "pkg.Func",
			err:        baseErr,
			wantString: "pkg.Func: base error",
			wantIs:     baseErr,
		},
		{
			name:       "wraps with op and message",
			op:         "pkg.Func",
			err:        baseErr,
			opts:       []Option{WithMsg("extra context")},
			wantString: "pkg.Func: base error: extra context",
			wantIs:     baseErr,
		},
		{
			name:       "wraps different error",
			op:         "other.Op",
			err:        otherErr,
			wantString: "other.Op: other error",
			wantIs:     otherErr,
		},
		{
			name:       "wraps with method style op",
			op:         "pkg.(Type).Method",
			err:        fmt.Errorf("something failed"),
			wantString: "pkg.(Type).Method: something failed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := Wrap(tt.op, tt.err, tt.opts...)
			require.NotNil(t, err)
			assert.Equal(t, tt.wantString, err.Error())
			if tt.wantIs != nil {
				assert.ErrorIs(t, err, tt.wantIs)
			}
		})
	}
}

func TestWrapNested(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		innerOp Op
		outerOp Op
		base    error
	}{
		{
			name:    "nested wrap preserves error chain",
			innerOp: "inner.Op",
			outerOp: "outer.Op",
			base:    errors.New("root cause"),
		},
		{
			name:    "triple nested wrap",
			innerOp: "a.Op",
			outerOp: "b.Op",
			base:    errors.New("deep error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			inner := Wrap(tt.innerOp, tt.base)
			outer := Wrap(tt.outerOp, inner)
			assert.ErrorIs(t, outer, tt.base)
		})
	}
}

func TestGetOpts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		opts    []Option
		wantMsg string
	}{
		{
			name:    "no options returns defaults",
			opts:    nil,
			wantMsg: "",
		},
		{
			name:    "single WithMsg",
			opts:    []Option{WithMsg("hello")},
			wantMsg: "hello",
		},
		{
			name:    "last WithMsg wins",
			opts:    []Option{WithMsg("first"), WithMsg("second")},
			wantMsg: "second",
		},
		{
			name:    "empty WithMsg",
			opts:    []Option{WithMsg("")},
			wantMsg: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			opts := GetOpts(tt.opts...)
			assert.Equal(t, tt.wantMsg, opts.withMsg)
		})
	}
}
