# Queue

`queue` is a generic, thread safe FIFO queue implementation in Go. It exposes a public `Queue[T]` without exposing the underlying implementation.

## Installation

```
go get github.com/RyanDerr/queue
```

## Usage

```go
package main

import (
	"fmt"

	"queue"
)

func main() {
	q := queue.New[string]()

	q.Enqueue("first")
	q.Enqueue("second")
	q.Enqueue("third")

	fmt.Println(q.Len())

	v, _ := q.Peak()
	fmt.Println(v)

	v, _ = q.Dequeue()
	fmt.Println(v)

	fmt.Println(q.Len())

	q.Clear()
	fmt.Println(q.Len())
}
```

## API

The `Queue[T]` interface provides the following methods

- `Len()` returns the number of elements in the queue
- `Enqueue(v T)` adds an element to the back of the queue
- `Dequeue()` removes and returns the element at the front of the queue
- `Peak()` returns the element at the front of the queue without removing it
- `Clear()` removes all elements from the queue

## Errors

Sentinel errors are provided for use with `errors.Is`

- `ErrEmptyQueue` is returned when an operation is attempted on an empty queue
- `ErrNilNode` is returned when an operation is performed on or with a nil node
- `ErrInitialization` is returned when the queue fails to initialize properly
- `ErrInternal` is returned when unexpected internal behavior is encountered
