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
func (sf *QuickQueue[T]) Len() int { return len(sf.head) - sf.headPos + len(sf.tail) }

// IsEmpty returns true if this Queue contains no elements.
func (sf *QuickQueue[T]) IsEmpty() bool { return sf.Len() == 0 }

// Clear initializes or clears queue.
func (sf *QuickQueue[T]) Clear() { sf.head, sf.tail, sf.headPos = nil, nil, 0 } // should set nil for gc

// Add items to the queue.
func (sf *QuickQueue[T]) Add(v T) { sf.tail = append(sf.tail, v) }

// Peek retrieves, but does not remove, the head of this Queue, or return nil if this Queue is empty.
func (sf *QuickQueue[T]) Peek() (v T, ok bool) {
	if sf.headPos < len(sf.head) {
		return sf.head[sf.headPos], true
	}
	if len(sf.tail) > 0 {
		return sf.tail[0], true
	}
	return v, false
}

// Poll retrieves and removes the head of the this Queue, or return nil if this Queue is empty.
func (sf *QuickQueue[T]) Poll() (v T, ok bool) {
	var placeholder T

	if sf.headPos >= len(sf.head) {
		if len(sf.tail) == 0 {
			return v, false
		}
		// Pick up tail as new head, clear tail.
		sf.head, sf.headPos, sf.tail = sf.tail, 0, sf.head[:0]
	}
	v = sf.head[sf.headPos]
	sf.head[sf.headPos] = placeholder // should set nil for gc
	sf.headPos++
	return v, true
}

// Contains returns true if this queue contains the specified element.
func (sf *QuickQueue[T]) Contains(val T) bool {
	for i := sf.headPos; i < len(sf.head); i++ {
		if sf.head[i] == val {
			return true
		}
	}
	for _, v := range sf.tail {
		if v == val {
			return true
		}
	}
	return false
}

// Remove a single instance of the specified element from this queue, if it is present.
func (sf *QuickQueue[T]) Remove(val T) {
	var found bool
	var idx int
	var placeholder T

	for i := sf.headPos; i < len(sf.head); i++ {
		if sf.head[i] == val {
			idx, found = i, true
			break
		}
	}
	if found {
		if (idx - sf.headPos) < (len(sf.head)-sf.headPos)/2 {
			moveLastToFirst(sf.head[sf.headPos:(idx + 1)])
			sf.head[sf.headPos] = placeholder // should set nil for gc
			sf.headPos++
		} else {
			moveFirstToLast(sf.head[idx:])
			sf.head[len(sf.head)-1] = placeholder // should set nil for gc
			sf.head = sf.head[:len(sf.head)-1]
		}
		return
	}

	for i, v := range sf.tail {
		if v == val {
			idx, found = i, true
			break
		}
	}
	if found {
		moveFirstToLast(sf.tail[idx:])
		sf.tail[len(sf.tail)-1] = placeholder // should set nil for gc
		sf.tail = sf.tail[:len(sf.tail)-1]
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
