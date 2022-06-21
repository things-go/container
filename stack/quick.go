package stack

// var _ container.Stack = (*QuickStack[string])(nil)

// QuickStack is quick LIFO stack implement with slice.
type QuickStack[T any] struct {
	items []T
}

// NewQuickStack creates a QuickStack. which implement interface stack.Interface.
func NewQuickStack[T any]() *QuickStack[T] { return &QuickStack[T]{} }

// Len returns the length of this priority queue.
func (sf *QuickStack[T]) Len() int { return len(sf.items) }

// IsEmpty returns true if this QuickStack contains no elements.
func (sf *QuickStack[T]) IsEmpty() bool { return len(sf.items) == 0 }

// Clear removes all the elements from this QuickStack.
func (sf *QuickStack[T]) Clear() { sf.items = nil } // should set nil for gc

// Push an element into this QuickStack.
func (sf *QuickStack[T]) Push(val T) { sf.items = append(sf.items, val) }

// Pop the element on the top of this QuickStack.
// return nil if this QuickStack is empty.
func (sf *QuickStack[T]) Pop() (v T, ok bool) {
	if length := len(sf.items); length > 0 {
		val := sf.items[length-1]
		sf.items = sf.items[:length-1]
		return val, true
	}
	return v, false
}

// Peek retrieves, but does not remove,
// the element on the top of this QuickStack,
// or return nil if this QuickStack is empty.
func (sf *QuickStack[T]) Peek() (v T, ok bool) {
	if len(sf.items) > 0 {
		return sf.items[len(sf.items)-1], true
	}
	return v, false
}

// Copy returns a copy of this stack.
func (sf *QuickStack[T]) Clone() *QuickStack[T] {
	items := make([]T, len(sf.items))
	copy(items, sf.items)
	return &QuickStack[T]{items}
}
