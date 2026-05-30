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

	"github.com/RyanDerr/queue"
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
- `DequeueLeft()` removes and returns the element at the back of the queue
- `Peak()` returns the element at the front of the queue without removing it
- `Clear()` removes all elements from the queue

## Options

`New` accepts functional options to configure the queue at construction time. If no options are provided the queue is unbounded.

| Option                 | Description                                                                                                                                     |
| ---------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| `WithCapacity(n uint)` | Sets a maximum number of elements the queue can hold. `Enqueue` returns `ErrQueueFull` when the limit is reached. A value of 0 means unlimited. |

```go
// Bounded queue that holds at most 100 elements
q := queue.New[string](queue.WithCapacity(100))

err := q.Enqueue("item")
if errors.Is(err, queue.ErrQueueFull) {
    // handle full queue
}
```

## Errors

Sentinel errors are provided for use with `errors.Is`

- `ErrEmptyQueue` is returned when an operation is attempted on an empty queue
- `ErrNilNode` is returned when an operation is performed on or with a nil node
- `ErrQueueFull` is returned when `Enqueue` is called on a queue that has reached its capacity
- `ErrInitialization` is returned when the queue fails to initialize properly
- `ErrInternal` is returned when unexpected internal behavior is encountered
