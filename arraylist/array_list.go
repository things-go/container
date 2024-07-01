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

// Package arraylist implements both an array List.
package arraylist

import (
	"fmt"
	"slices"

	"github.com/things-go/container"
)

var _ container.List[int] = (*List[int])(nil)

// List represents an array list.
// It implements the interface list.Interface.
type List[T comparable] struct {
	items []T
}

// New initializes and returns an ArrayList.
func New[T comparable]() *List[T] {
	return &List[T]{items: []T{}}
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *List[T]) Len() int { return len(l.items) }

// IsEmpty returns the list l is empty or not.
func (l *List[T]) IsEmpty() bool { return l.Len() == 0 }

// Clear initializes or clears list l.
func (l *List[T]) Clear() { l.items = make([]T, 0) }

// Push inserts a new element e with value v at the back of list l.
func (l *List[T]) Push(items T) { l.items = append(l.items, items) }

// PushFront inserts a new element e with value v at the front of list l.
func (l *List[T]) PushFront(v T) {
	l.items = append(l.items, v)
	moveLastToFirst(l.items)
}

// PushBack inserts a new element e with value v at the back of list l.
func (l *List[T]) PushBack(v T) { l.items = append(l.items, v) }

// Add inserts the specified element at the specified position in this list.
func (l *List[T]) Add(index int, val T) error {
	if index < 0 || index > len(l.items) {
		return fmt.Errorf("index out of range, index: %d, len: %d", index, l.Len())
	}

	if index == l.Len() {
		l.Push(val)
	} else {
		length := len(l.items)
		l.items = append(l.items, val)
		copy(l.items[index+1:], l.items[index:length])
		l.items[index] = val
	}
	return nil
}

// PushFrontList inserts a copy of an other list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (l *List[T]) PushFrontList(other *List[T]) {
	items := make([]T, 0, len(l.items)+len(other.items))
	items = append(items, other.items...)
	items = append(items, l.items...)
	l.items = items
}

// PushBackList inserts a copy of an other list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *List[T]) PushBackList(other *List[T]) {
	l.items = append(l.items, other.items...)
}

// Poll return the front element value and then remove from list.
func (l *List[T]) Poll() (val T, ok bool) {
	return l.PollFront()
}

// PollFront return the front element value and then remove from list.
func (l *List[T]) PollFront() (val T, ok bool) {
	var placeholder T

	if n := len(l.items); n > 0 {
		moveFirstToLast(l.items)
		val = l.items[n-1]
		l.items[n-1] = placeholder // for gc
		l.items = l.items[:n-1]
		ok = true
	}
	return val, ok
}

// PollBack return the back element value and then remove from list.
func (l *List[T]) PollBack() (val T, ok bool) {
	var placeholder T

	if n := len(l.items); n > 0 {
		val = l.items[n-1]
		l.items[n-1] = placeholder // for gc
		l.items = l.items[:n-1]
		ok = true
	}
	return val, ok
}

// Remove removes the element at the specified position in this list.
// It returns an error if the index is out of range.
func (l *List[T]) Remove(index int) (val T, err error) {
	var placeholder T

	if index < 0 || index >= len(l.items) {
		return val, fmt.Errorf("index out of range, index: %d, len: %d", index, l.Len())
	}

	val = l.items[index]
	moveFirstToLast(l.items[index:])
	l.items[len(l.items)-1] = placeholder
	l.items = l.items[:len(l.items)-1]
	l.shrinkList()
	return val, nil
}

// RemoveValue removes the first occurrence of the specified element from this list, if it is present.
// It returns false if the target value isn't present, otherwise returns true.
func (l *List[T]) RemoveValue(val T) bool {
	var placeholder T

	if l.Len() == 0 {
		return false
	}

	if idx := l.indexOf(val); idx >= 0 {
		moveFirstToLast(l.items[idx:])
		l.items[len(l.items)-1] = placeholder
		l.items = l.items[:len(l.items)-1]
		l.shrinkList()
		return true
	}
	return false
}

// Get returns the element at the specified position in this list. The index must be in the range of [0, size).
func (l *List[T]) Get(index int) (val T, err error) {
	if index < 0 || index >= len(l.items) {
		return val, fmt.Errorf("index out of range, index:%d, len:%d", index, l.Len())
	}

	return l.items[index], nil
}

// Peek return the front element value.
func (l *List[T]) Peek() (val T, ok bool) {
	return l.PeekFront()
}

// PeekFront return the front element value.
func (l *List[T]) PeekFront() (val T, ok bool) {
	if len(l.items) > 0 {
		return l.items[0], true
	}
	return val, false
}

// PeekBack return the back element value.
func (l *List[T]) PeekBack() (val T, ok bool) {
	if len(l.items) > 0 {
		return l.items[len(l.items)-1], true
	}
	return val, false
}

// Iterator returns an iterator over the elements in this list in proper sequence.
func (l *List[T]) Iterator(f func(T) bool) {
	for index := 0; index < l.Len(); index++ {
		if f == nil || !f(l.items[index]) {
			return
		}
	}
}

// ReverseIterator returns an iterator over the elements in this list in reverse sequence as Iterator.
func (l *List[T]) ReverseIterator(f func(T) bool) {
	for index := l.Len() - 1; index >= 0; index-- {
		if f == nil || !f(l.items[index]) {
			return
		}
	}
}

// Contains contains the value.
func (l *List[T]) Contains(val T) bool {
	return l.indexOf(val) >= 0
}

// Sort the list.
func (l *List[T]) Sort(less func(a, b T) int) {
	slices.SortFunc(l.items, less)
}

// Values get a copy of all the values in the list.
func (l *List[T]) Values() []T {
	return slices.Clone(l.items)
}

func (l *List[T]) shrinkList() {
	oldLen, oldCap := len(l.items), cap(l.items)
	if oldCap > 1024 && oldLen <= oldCap/4 { // shrink when len(list) <= cap(list)/4
		newItems := make([]T, oldLen)
		copy(newItems, l.items)
		l.Clear()
		l.items = newItems
	}
}

// indexOf returns the index of the first occurrence of the specified element
// in this list, or -1 if this list does not contain the element.
func (l *List[T]) indexOf(val T) int {
	for i, v := range l.items {
		if v == val {
			return i
		}
	}
	return -1
}

// move the last element to the first position.
// for example: [1,2,3,4] -> [4,1,2,3]
func moveLastToFirst[T any](items []T) {
	for i := 0; i < len(items)-1; i++ {
		items[i], items[len(items)-1] = items[len(items)-1], items[i]
	}
}

// move the first element to the last position.
// for example: [1,2,3,4] -> [2,3,4,1]
func moveFirstToLast[T any](items []T) {
	for i := 0; i < len(items)-1; i++ {
		items[0], items[len(items)-1-i] = items[len(items)-1-i], items[0]
	}
}
