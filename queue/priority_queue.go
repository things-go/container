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
	"container/heap"

	"github.com/things-go/container"
)

var _ container.Queue[container.Int] = (*PriorityQueue[container.Int])(nil)

// PriorityQueue represents an unbounded priority queue based on a priority heap.
// It implements heap.Interface.
type PriorityQueue[T container.Comparable] struct {
	data *container.Container[T]
}

// NewPriorityQueue initializes and returns an Queue, default min heap.
func NewPriorityQueue[T container.Comparable](maxHeap bool, items ...T) *PriorityQueue[T] {
	q := &PriorityQueue[T]{
		data: &container.Container[T]{
			Items:   items,
			Reverse: maxHeap,
		},
	}
	heap.Init(q.data)
	return q
}

// Len returns the length of this priority queue.
func (pq *PriorityQueue[T]) Len() int { return pq.data.Len() }

// IsEmpty returns true if this list contains no elements.
func (pq *PriorityQueue[T]) IsEmpty() bool { return pq.Len() == 0 }

// Clear removes all the elements from this priority queue.
func (pq *PriorityQueue[T]) Clear() { pq.data.Items = make([]T, 0) }

// Add inserts the specified element into this priority queue.
func (pq *PriorityQueue[T]) Add(items T) {
	heap.Push(pq.data, items)
}

// Peek retrieves, but does not remove, the head of this queue, or return nil if this queue is empty.
func (pq *PriorityQueue[T]) Peek() (val T, exist bool) {
	if pq.Len() > 0 {
		return pq.data.Items[0], true
	}
	return val, false
}

// Poll retrieves and removes the head of the this queue, or return nil if this queue is empty.
func (pq *PriorityQueue[T]) Poll() (val T, exist bool) {
	if pq.Len() > 0 {
		return heap.Pop(pq.data).(T), true
	}
	return val, false
}

// Contains returns true if this queue contains the specified element.
func (pq *PriorityQueue[T]) Contains(val T) bool { return pq.indexOf(val) >= 0 }

// Remove a single instance of the specified element from this queue, if it is present.
// It returns false if the target value isn't present, otherwise returns true.
func (pq *PriorityQueue[T]) Remove(val T) {
	if idx := pq.indexOf(val); idx >= 0 {
		heap.Remove(pq.data, idx)
	}
}

func (pq *PriorityQueue[T]) indexOf(val T) int {
	if pq.Len() > 0 {
		for i := 0; i < pq.Len(); i++ {
			if val.CompareTo(pq.data.Items[i]) == 0 {
				return i
			}
		}
	}
	return -1
}
