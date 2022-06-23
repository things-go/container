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
	"golang.org/x/exp/constraints"

	"github.com/things-go/container"
	"github.com/things-go/container/core/heap"
)

var _ container.Queue[int] = (*PriorityQueue[int])(nil)

// PriorityQueue represents an unbounded priority queue based on a priority heap.
// It implements heap.Interface.
type PriorityQueue[T constraints.Ordered] struct {
	data *container.Container[T]
}

// Option for NewPriorityQueue.
type Option[T constraints.Ordered] func(q *PriorityQueue[T])

// NewPriorityQueue initializes and returns an Queue, default min heap.
func NewPriorityQueue[T constraints.Ordered](maxHeap bool, items ...T) *PriorityQueue[T] {
	q := &PriorityQueue[T]{
		data: &container.Container[T]{
			items,
			maxHeap,
		},
	}
	heap.Init[T](q.data)
	return q
}

// Len returns the length of this priority queue.
func (sf *PriorityQueue[T]) Len() int { return sf.data.Len() }

// IsEmpty returns true if this list contains no elements.
func (sf *PriorityQueue[T]) IsEmpty() bool { return sf.Len() == 0 }

// Clear removes all the elements from this priority queue.
func (sf *PriorityQueue[T]) Clear() { sf.data.Items = make([]T, 0) }

// Add inserts the specified element into this priority queue.
func (sf *PriorityQueue[T]) Add(items T) {
	heap.Push[T](sf.data, items)
}

// Peek retrieves, but does not remove, the head of this queue, or return nil if this queue is empty.
func (sf *PriorityQueue[T]) Peek() (val T, exist bool) {
	if sf.Len() > 0 {
		return sf.data.Items[0], true
	}
	return val, false
}

// Poll retrieves and removes the head of the this queue, or return nil if this queue is empty.
func (sf *PriorityQueue[T]) Poll() (val T, exist bool) {
	if sf.Len() > 0 {
		return heap.Pop[T](sf.data), true
	}
	return val, false
}

// Contains returns true if this queue contains the specified element.
func (sf *PriorityQueue[T]) Contains(val T) bool { return sf.indexOf(val) >= 0 }

// Remove a single instance of the specified element from this queue, if it is present.
// It returns false if the target value isn't present, otherwise returns true.
func (sf *PriorityQueue[T]) Remove(val T) {
	if idx := sf.indexOf(val); idx >= 0 {
		heap.Remove[T](sf.data, idx)
	}
}

func (sf *PriorityQueue[T]) indexOf(val T) int {
	if sf.Len() > 0 {
		for i := 0; i < sf.Len(); i++ {
			if val == sf.data.Items[i] {
				return i
			}
		}
	}
	return -1
}
