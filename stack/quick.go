package stack

import (
	"github.com/things-go/container"
)

var _ container.Stack[string] = (*QuickStack[string])(nil)

// QuickStack is quick LIFO stack implement with slice.
type QuickStack[T any] struct {
	items []T
}

// NewQuickStack creates a QuickStack. which implement interface stack.Interface.
func NewQuickStack[T any]() *QuickStack[T] { return &QuickStack[T]{} }

// Len returns the length of this priority queue.
func (qs *QuickStack[T]) Len() int { return len(qs.items) }

// IsEmpty returns true if this QuickStack contains no elements.
func (qs *QuickStack[T]) IsEmpty() bool { return len(qs.items) == 0 }

// Clear removes all the elements from this QuickStack.
func (qs *QuickStack[T]) Clear() { qs.items = nil } // should set nil for gc

// Push an element into this QuickStack.
func (qs *QuickStack[T]) Push(val T) { qs.items = append(qs.items, val) }

// Pop the element on the top of this QuickStack.
// return nil if this QuickStack is empty.
func (qs *QuickStack[T]) Pop() (v T, ok bool) {
	if length := len(qs.items); length > 0 {
		val := qs.items[length-1]
		qs.items = qs.items[:length-1]
		return val, true
	}
	return v, false
}

// Peek retrieves, but does not remove,
// the element on the top of this QuickStack,
// or return nil if this QuickStack is empty.
func (qs *QuickStack[T]) Peek() (v T, ok bool) {
	if len(qs.items) > 0 {
		return qs.items[len(qs.items)-1], true
	}
	return v, false
}

// Copy returns a copy of this stack.
func (qs *QuickStack[T]) Clone() *QuickStack[T] {
	items := make([]T, len(qs.items))
	copy(items, qs.items)
	return &QuickStack[T]{items}
}
