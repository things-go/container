// Copyright [2022] [thinkgos]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package queue

import (
	"github.com/things-go/container"
)

var _ container.Queue[int] = (*Queue[int])(nil)

// element is an element of the Queue implement with list.
type element[T comparable] struct {
	next  *element[T]
	value T
}

// Queue represents a singly linked list.
type Queue[T comparable] struct {
	head   *element[T]
	tail   *element[T]
	length int
}

// New creates a Queue. which implement queue.Interface.
func New[T comparable]() *Queue[T] {
	return new(Queue[T])
}

// Len returns the length of this queue.
func (sf *Queue[T]) Len() int { return sf.length }

// IsEmpty returns true if this Queue contains no elements.
func (sf *Queue[T]) IsEmpty() bool { return sf.Len() == 0 }

// Clear initializes or clears queue.
func (sf *Queue[T]) Clear() { sf.head, sf.tail, sf.length = nil, nil, 0 }

// Add items to the queue.
func (sf *Queue[T]) Add(v T) {
	e := &element[T]{value: v}
	if sf.tail == nil {
		sf.head, sf.tail = e, e
	} else {
		sf.tail.next = e
		sf.tail = e
	}
	sf.length++
}

// Peek retrieves, but does not remove, the head of this Queue, or return nil if this Queue is empty.
func (sf *Queue[T]) Peek() (v T, ok bool) {
	if sf.head != nil {
		return sf.head.value, true
	}
	return v, false
}

// Poll retrieves and removes the head of the this Queue, or return nil if this Queue is empty.
func (sf *Queue[T]) Poll() (v T, ok bool) {
	if sf.head != nil {
		v = sf.head.value
		sf.head = sf.head.next
		if sf.head == nil {
			sf.tail = nil
		}
		sf.length--
		ok = true
	}
	return v, ok
}

// Contains returns true if this queue contains the specified element.
func (sf *Queue[T]) Contains(val T) bool {
	for e := sf.head; e != nil; e = e.next {
		if val == e.value {
			return true
		}
	}
	return false
}

// Remove a single instance of the specified element from this queue, if it is present.
func (sf *Queue[T]) Remove(val T) {
	for pre, e := sf.head, sf.head; e != nil; {
		if val == e.value {
			if sf.head == e && sf.tail == e {
				sf.head, sf.tail = nil, nil
			} else if sf.head == e {
				sf.head = e.next
			} else if sf.tail == e {
				sf.tail = pre
				sf.tail.next = nil
			} else {
				pre.next = e.next
			}
			e.next = nil
			sf.length--
			return
		}
		pre = e
		e = e.next
	}
}
