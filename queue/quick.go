package queue

import (
	"github.com/things-go/container"
)

var _ container.Queue[int] = (*QuickQueue[int])(nil)

// QuickQueue implement with slice.
type QuickQueue[T comparable] struct {
	headPos int
	head    []T
	tail    []T
}

// NewQuickQueue new quick queue.
func NewQuickQueue[T comparable]() *QuickQueue[T] {
	return new(QuickQueue[T])
}

// Len returns the length of this queue.
func (q *QuickQueue[T]) Len() int { return len(q.head) - q.headPos + len(q.tail) }

// IsEmpty returns true if this Queue contains no elements.
func (q *QuickQueue[T]) IsEmpty() bool { return q.Len() == 0 }

// Clear initializes or clears queue.
func (q *QuickQueue[T]) Clear() { q.head, q.tail, q.headPos = nil, nil, 0 } // should set nil for gc

// Add items to the queue.
func (q *QuickQueue[T]) Add(v T) { q.tail = append(q.tail, v) }

// Peek retrieves, but does not remove, the head of this Queue, or return nil if this Queue is empty.
func (q *QuickQueue[T]) Peek() (v T, ok bool) {
	if q.headPos < len(q.head) {
		return q.head[q.headPos], true
	}
	if len(q.tail) > 0 {
		return q.tail[0], true
	}
	return v, false
}

// Poll retrieves and removes the head of the this Queue, or return nil if this Queue is empty.
func (q *QuickQueue[T]) Poll() (v T, ok bool) {
	var placeholder T

	if q.headPos >= len(q.head) {
		if len(q.tail) == 0 {
			return v, false
		}
		// Pick up tail as new head, clear tail.
		q.head, q.headPos, q.tail = q.tail, 0, q.head[:0]
	}
	v = q.head[q.headPos]
	q.head[q.headPos] = placeholder // should set nil for gc
	q.headPos++
	return v, true
}

// Contains returns true if this queue contains the specified element.
func (q *QuickQueue[T]) Contains(val T) bool {
	for i := q.headPos; i < len(q.head); i++ {
		if q.head[i] == val {
			return true
		}
	}
	for _, v := range q.tail {
		if v == val {
			return true
		}
	}
	return false
}

// Remove a single instance of the specified element from this queue, if it is present.
func (q *QuickQueue[T]) Remove(val T) {
	var found bool
	var idx int
	var placeholder T

	for i := q.headPos; i < len(q.head); i++ {
		if q.head[i] == val {
			idx, found = i, true
			break
		}
	}
	if found {
		if (idx - q.headPos) < (len(q.head)-q.headPos)/2 {
			moveLastToFirst(q.head[q.headPos:(idx + 1)])
			q.head[q.headPos] = placeholder // should set nil for gc
			q.headPos++
		} else {
			moveFirstToLast(q.head[idx:])
			q.head[len(q.head)-1] = placeholder // should set nil for gc
			q.head = q.head[:len(q.head)-1]
		}
		return
	}

	for i, v := range q.tail {
		if v == val {
			idx, found = i, true
			break
		}
	}
	if found {
		moveFirstToLast(q.tail[idx:])
		q.tail[len(q.tail)-1] = placeholder // should set nil for gc
		q.tail = q.tail[:len(q.tail)-1]
	}
}

func moveLastToFirst[T any](items []T) {
	for i := 0; i < len(items); i++ {
		items[i], items[len(items)-1] = items[len(items)-1], items[i]
	}
}

func moveFirstToLast[T any](items []T) {
	for i := 0; i < len(items); i++ {
		items[0], items[len(items)-1-i] = items[len(items)-1-i], items[0]
	}
}
