package ringbuffer

import "sync"

type RingBuffer[T any] struct {
	mu   sync.RWMutex
	buf  []T
	wptr int // write pointer: next index to write
	size int // current number of elements
	cap  int // maximum capacity (len(buf))
}

// New creates a new RingBuffer with the given capacity (must be > 0).
func New[T any](capacity int) *RingBuffer[T] {
	if capacity <= 0 {
		capacity = 128
	}
	return &RingBuffer[T]{
		buf:  make([]T, capacity),
		wptr: 0,
		size: 0,
		cap:  capacity,
	}
}

// Push adds a value to the buffer.
// If the buffer is full, the oldest element is overwritten.
func (rb *RingBuffer[T]) Push(value T) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.buf[rb.wptr] = value
	rb.wptr = (rb.wptr + 1) % rb.cap

	if rb.size < rb.cap {
		rb.size++
	}
}

// Len returns the current number of elements in the buffer.
func (rb *RingBuffer[T]) Len() int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.size
}

// Cap returns the maximum capacity of the buffer.
func (rb *RingBuffer[T]) Cap() int {
	return rb.cap
}

// Tail returns the oldest up to n elements in insertion order (oldest first).
// Returns at most min(n, Len()) elements.
func (rb *RingBuffer[T]) Tail(n int) []T {
	if n <= 0 {
		return nil
	}

	rb.mu.RLock()
	defer rb.mu.RUnlock()

	if rb.size == 0 {
		return []T{}
	}

	count := n
	if count > rb.size {
		count = rb.size
	}

	result := make([]T, count)
	// Oldest element is at (wptr - size + cap) % cap
	readStart := (rb.wptr - rb.size + rb.cap) % rb.cap
	for i := 0; i < count; i++ {
		idx := (readStart + i) % rb.cap
		result[i] = rb.buf[idx]
	}
	return result
}

// Head returns the newest up to n elements in insertion order (oldest among them first).
// For example, if buffer contains [A, B, C, D, E] (E latest), Head(3) returns [C, D, E].
// Returns at most min(n, Len()) elements.
func (rb *RingBuffer[T]) Head(n int) []T {
	if n <= 0 {
		return nil
	}

	rb.mu.RLock()
	defer rb.mu.RUnlock()

	if rb.size == 0 {
		return []T{}
	}

	count := n
	if count > rb.size {
		count = rb.size
	}

	result := make([]T, count)
	// Start of the "newest n" block:
	// It begins at (wptr - count + cap) % cap
	readStart := (rb.wptr - count + rb.cap) % rb.cap
	for i := 0; i < count; i++ {
		idx := (readStart + i) % rb.cap
		result[i] = rb.buf[idx]
	}
	return result
}
