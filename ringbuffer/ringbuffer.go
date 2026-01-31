package ringbuffer

import "sync"

type RingBuffer[T any] struct {
	mu   sync.RWMutex
	buf  []T
	wptr int
	size int
	cap  int
}

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

func (rb *RingBuffer[T]) Push(value T) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.buf[rb.wptr] = value
	rb.wptr = (rb.wptr + 1) % rb.cap

	if rb.size < rb.cap {
		rb.size++
	}
}

func (rb *RingBuffer[T]) Len() int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.size
}

func (rb *RingBuffer[T]) Cap() int {
	return rb.cap
}

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
	readStart := (rb.wptr - rb.size + rb.cap) % rb.cap
	for i := 0; i < count; i++ {
		idx := (readStart + i) % rb.cap
		result[i] = rb.buf[idx]
	}
	return result
}

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
	readStart := (rb.wptr - count + rb.cap) % rb.cap
	for i := 0; i < count; i++ {
		idx := (readStart + i) % rb.cap
		result[i] = rb.buf[idx]
	}
	return result
}
