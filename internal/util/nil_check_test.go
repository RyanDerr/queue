package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNil(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input any
		want  bool
	}{
		{
			name:  "nil interface",
			input: nil,
			want:  true,
		},
		{
			name:  "nil int pointer",
			input: (*int)(nil),
			want:  true,
		},
		{
			name:  "nil string pointer",
			input: (*string)(nil),
			want:  true,
		},
		{
			name:  "nil struct pointer",
			input: (*struct{})(nil),
			want:  true,
		},
		{
			name:  "non-nil int pointer",
			input: ptrTo(42),
			want:  false,
		},
		{
			name:  "non-nil string pointer",
			input: ptrTo("hello"),
			want:  false,
		},
		{
			name:  "nil slice",
			input: ([]int)(nil),
			want:  true,
		},
		{
			name:  "empty slice",
			input: []int{},
			want:  false,
		},
		{
			name:  "nil string slice",
			input: ([]string)(nil),
			want:  true,
		},
		{
			name:  "empty string slice",
			input: []string{},
			want:  false,
		},
		{
			name:  "nil map",
			input: (map[string]int)(nil),
			want:  true,
		},
		{
			name:  "empty map",
			input: map[string]int{},
			want:  false,
		},
		{
			name:  "nil chan int",
			input: (chan int)(nil),
			want:  true,
		},
		{
			name:  "nil chan string",
			input: (chan string)(nil),
			want:  true,
		},
		{
			name:  "nil func",
			input: (func())(nil),
			want:  true,
		},
		{
			name:  "int zero value",
			input: 0,
			want:  false,
		},
		{
			name:  "int non-zero",
			input: 42,
			want:  false,
		},
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
		{
			name:  "non-empty string",
			input: "test",
			want:  false,
		},
		{
			name:  "false bool",
			input: false,
			want:  false,
		},
		{
			name:  "true bool",
			input: true,
			want:  false,
		},
		{
			name:  "float64 zero",
			input: 0.0,
			want:  false,
		},
		{
			name:  "float64 non-zero",
			input: 3.14,
			want:  false,
		},
		{
			name:  "empty struct",
			input: struct{}{},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, IsNil(tt.input))
		})
	}
}

func ptrTo[T any](v T) *T {
	return &v
}
