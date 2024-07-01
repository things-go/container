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
func (q *Queue[T]) Len() int { return q.length }

// IsEmpty returns true if this Queue contains no elements.
func (q *Queue[T]) IsEmpty() bool { return q.Len() == 0 }

// Clear initializes or clears queue.
func (q *Queue[T]) Clear() { q.head, q.tail, q.length = nil, nil, 0 }

// Add items to the queue.
func (q *Queue[T]) Add(v T) {
	e := &element[T]{value: v}
	if q.tail == nil {
		q.head, q.tail = e, e
	} else {
		q.tail.next = e
		q.tail = e
	}
	q.length++
}

// Peek retrieves, but does not remove, the head of this Queue, or return nil if this Queue is empty.
func (q *Queue[T]) Peek() (v T, ok bool) {
	if q.head != nil {
		return q.head.value, true
	}
	return v, false
}

// Poll retrieves and removes the head of the this Queue, or return nil if this Queue is empty.
func (q *Queue[T]) Poll() (v T, ok bool) {
	if q.head != nil {
		v = q.head.value
		q.head = q.head.next
		if q.head == nil {
			q.tail = nil
		}
		q.length--
		ok = true
	}
	return v, ok
}

// Contains returns true if this queue contains the specified element.
func (q *Queue[T]) Contains(val T) bool {
	for e := q.head; e != nil; e = e.next {
		if val == e.value {
			return true
		}
	}
	return false
}

// Remove a single instance of the specified element from this queue, if it is present.
func (q *Queue[T]) Remove(val T) {
	for pre, e := q.head, q.head; e != nil; {
		if val == e.value {
			if q.head == e && q.tail == e {
				q.head, q.tail = nil, nil
			} else if q.head == e {
				q.head = e.next
			} else if q.tail == e {
				q.tail = pre
				q.tail.next = nil
			} else {
				pre.next = e.next
			}
			e.next = nil
			q.length--
			return
		}
		pre = e
		e = e.next
	}
}
