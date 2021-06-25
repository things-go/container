package stack

import (
	"github.com/things-go/container"
)

var _ container.Stack = (*QuickStack)(nil)

// QuickStack is quick LIFO stack implement with slice.
type QuickStack struct {
	items []interface{}
}

// NewQuickStack creates a QuickStack. which implement interface stack.Interface.
func NewQuickStack() *QuickStack { return &QuickStack{} }

// Len returns the length of this priority queue.
func (sf *QuickStack) Len() int { return len(sf.items) }

// IsEmpty returns true if this QuickStack contains no elements.
func (sf *QuickStack) IsEmpty() bool { return len(sf.items) == 0 }

// Clear removes all the elements from this QuickStack.
func (sf *QuickStack) Clear() { sf.items = nil } // should set nil for gc

// Push push an element into this QuickStack.
func (sf *QuickStack) Push(val interface{}) { sf.items = append(sf.items, val) }

// Pop pop the element on the top of this QuickStack.
// return nil if this QuickStack is empty.
func (sf *QuickStack) Pop() interface{} {
	if length := len(sf.items); length > 0 {
		val := sf.items[length-1]
		sf.items[length-1] = nil // should set nil for gc
		sf.items = sf.items[:length-1]
		return val
	}
	return nil
}

// Peek retrieves, but does not remove,
// the element on the top of this QuickStack,
// or return nil if this QuickStack is empty.
func (sf *QuickStack) Peek() interface{} {
	if len(sf.items) > 0 {
		return sf.items[len(sf.items)-1]
	}
	return nil
}
