package node

import "errors"

var (
	// ErrNilNode is returned when an operation is performed on or with a nil node.
	ErrNilNode = errors.New("node is nil")
)
