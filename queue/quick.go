package queue

import (
	"reflect"

	"github.com/things-go/container"
)

var _ container.Queue = (*QuickQueue)(nil)

// QuickQueue implement with slice.
type QuickQueue struct {
	headPos int
	head    []interface{}
	tail    []interface{}
	compare container.Comparator
}

// NewQuickQueue new quick queue.
func NewQuickQueue(opts ...Option) *QuickQueue {
	q := new(QuickQueue)
	for _, opt := range opts {
		opt(q)
	}
	return q
}

func (sf *QuickQueue) apply(c container.Comparator) { sf.compare = c }

// Len returns the length of this queue.
func (sf *QuickQueue) Len() int { return len(sf.head) - sf.headPos + len(sf.tail) }

// IsEmpty returns true if this Queue contains no elements.
func (sf *QuickQueue) IsEmpty() bool { return sf.Len() == 0 }

// Clear initializes or clears queue.
func (sf *QuickQueue) Clear() { sf.head, sf.tail, sf.headPos = nil, nil, 0 } // should set nil for gc

// Add items to the queue.
func (sf *QuickQueue) Add(v interface{}) { sf.tail = append(sf.tail, v) }

// Peek retrieves, but does not remove, the head of this Queue, or return nil if this Queue is empty.
func (sf *QuickQueue) Peek() interface{} {
	if sf.headPos < len(sf.head) {
		return sf.head[sf.headPos]
	}
	if len(sf.tail) > 0 {
		return sf.tail[0]
	}
	return nil
}

// Poll retrieves and removes the head of the this Queue, or return nil if this Queue is empty.
func (sf *QuickQueue) Poll() interface{} {
	if sf.headPos >= len(sf.head) {
		if len(sf.tail) == 0 {
			return nil
		}
		// Pick up tail as new head, clear tail.
		sf.head, sf.headPos, sf.tail = sf.tail, 0, sf.head[:0]
	}
	v := sf.head[sf.headPos]
	sf.head[sf.headPos] = nil // should set nil for gc
	sf.headPos++
	return v
}

// Contains returns true if this queue contains the specified element.
func (sf *QuickQueue) Contains(val interface{}) bool {
	for i := sf.headPos; i < len(sf.head); i++ {
		if sf.Compare(sf.head[i], val) {
			return true
		}
	}
	for _, v := range sf.tail {
		if sf.Compare(v, val) {
			return true
		}
	}
	return false
}

// Remove a single instance of the specified element from this queue, if it is present.
func (sf *QuickQueue) Remove(val interface{}) {
	var found bool
	var idx int

	for i := sf.headPos; i < len(sf.head); i++ {
		if sf.Compare(sf.head[i], val) {
			idx, found = i, true
			break
		}
	}
	if found {
		if (idx - sf.headPos) < (len(sf.head)-sf.headPos)/2 {
			moveLastToFirst(sf.head[sf.headPos:(idx + 1)])
			sf.head[sf.headPos] = nil // should set nil for gc
			sf.headPos++
		} else {
			moveFirstToLast(sf.head[idx:])
			sf.head[len(sf.head)-1] = nil // should set nil for gc
			sf.head = sf.head[:len(sf.head)-1]
		}
		return
	}

	for i, v := range sf.tail {
		if sf.Compare(v, val) {
			idx, found = i, true
			break
		}
	}
	if found {
		moveFirstToLast(sf.tail[idx:])
		sf.tail[len(sf.tail)-1] = nil // should set nil for gc
		sf.tail = sf.tail[:len(sf.tail)-1]
	}
}

func (sf *QuickQueue) Compare(v1, v2 interface{}) bool {
	if sf.compare == nil {
		return reflect.DeepEqual(v1, v2)
	}
	return sf.compare(v1, v2) == 0
}

func moveLastToFirst(items []interface{}) {
	for i := 0; i < len(items); i++ {
		items[i], items[len(items)-1] = items[len(items)-1], items[i]
	}
}

func moveFirstToLast(items []interface{}) {
	for i := 0; i < len(items); i++ {
		items[0], items[len(items)-1-i] = items[len(items)-1-i], items[0]
	}
}
