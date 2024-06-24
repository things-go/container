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

// Package linkedlist implements both an linked and a list.
package linkedlist

import (
	"fmt"

	"slices"

	"github.com/things-go/container"
	"github.com/things-go/container/core/list"
)

var _ container.List[int] = (*LinkedList[int])(nil)

// LinkedList represents a doubly linked list.
// It implements the interface list.Interface.
type LinkedList[T comparable] struct {
	l *list.List[T]
}

// New initializes and returns an LinkedList.
func New[T comparable]() *LinkedList[T] {
	return &LinkedList[T]{l: list.New[T]()}
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (sf *LinkedList[T]) Len() int { return sf.l.Len() }

// IsEmpty returns the list l is empty or not.
func (sf *LinkedList[T]) IsEmpty() bool { return sf.l.Len() == 0 }

// Clear initializes or clears list l.
func (sf *LinkedList[T]) Clear() { sf.l.Init() }

// Push inserts a new element e with value v at the back of list l.
func (sf *LinkedList[T]) Push(v T) { sf.l.PushBack(v) }

// PushFront inserts a new element e with value v at the front of list l.
func (sf *LinkedList[T]) PushFront(v T) { sf.l.PushFront(v) }

// PushBack inserts a new element e with value v at the back of list l.
func (sf *LinkedList[T]) PushBack(v T) { sf.l.PushBack(v) }

// Add add to the index of the list with value.
func (sf *LinkedList[T]) Add(index int, val T) error {
	if index < 0 || index > sf.Len() {
		return fmt.Errorf("index out of range, index: %d, len: %d", index, sf.Len())
	}

	if index == sf.Len() {
		sf.l.PushBack(val)
	} else {
		sf.l.InsertBefore(val, sf.getElement(index))
	}
	return nil
}

// PushFrontList inserts a copy of another list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (sf *LinkedList[T]) PushFrontList(other *LinkedList[T]) {
	sf.l.PushFrontList(other.l)
}

// PushBackList inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (sf *LinkedList[T]) PushBackList(other *LinkedList[T]) {
	sf.l.PushBackList(other.l)
}

// Poll return the front element value and then remove from list.
func (sf *LinkedList[T]) Poll() (T, bool) {
	return sf.PollFront()
}

// PollFront return the front element value and then remove from list.
func (sf *LinkedList[T]) PollFront() (val T, ok bool) {
	e := sf.l.Front()
	if e != nil {
		return sf.l.Remove(e), true
	}
	return val, false
}

// PollBack return the back element value and then remove from list.
func (sf *LinkedList[T]) PollBack() (val T, ok bool) {
	e := sf.l.Back()
	if e != nil {
		return sf.l.Remove(e), true
	}
	return val, false
}

// Remove remove the index in the list.
func (sf *LinkedList[T]) Remove(index int) (val T, err error) {
	if index < 0 || index >= sf.Len() {
		return val, fmt.Errorf("index out of range, index:%d, len:%d", index, sf.Len())
	}
	return sf.l.Remove(sf.getElement(index)), nil
}

// RemoveValue remove the value in the list.
func (sf *LinkedList[T]) RemoveValue(val T) bool {
	if sf.Len() == 0 {
		return false
	}

	for e := sf.l.Front(); e != nil; e = e.Next() {
		if val == e.Value {
			sf.l.Remove(e)
			return true
		}
	}
	return false
}

// Get the index in the list.
func (sf *LinkedList[T]) Get(index int) (val T, err error) {
	if index < 0 || index >= sf.Len() {
		return val, fmt.Errorf("index out of range, index: %d, len: %d", index, sf.Len())
	}
	return sf.getElement(index).Value, nil
}

// Peek return the front element value.
func (sf *LinkedList[T]) Peek() (T, bool) {
	return sf.PeekFront()
}

// PeekFront return the front element value.
func (sf *LinkedList[T]) PeekFront() (val T, ok bool) {
	if e := sf.l.Front(); e != nil {
		return e.Value, true
	}
	return val, false
}

// PeekBack return the back element value.
func (sf *LinkedList[T]) PeekBack() (val T, ok bool) {
	if e := sf.l.Back(); e != nil {
		return e.Value, true
	}
	return val, false
}

// Iterator the list.
func (sf *LinkedList[T]) Iterator(cb func(T) bool) {
	for e := sf.l.Front(); e != nil; e = e.Next() {
		if cb == nil || !cb(e.Value) {
			return
		}
	}
}

// ReverseIterator reverse iterator the list.
func (sf *LinkedList[T]) ReverseIterator(cb func(T) bool) {
	for e := sf.l.Back(); e != nil; e = e.Prev() {
		if cb == nil || !cb(e.Value) {
			return
		}
	}
}

// Contains contains the value.
func (sf *LinkedList[T]) Contains(val T) bool {
	return sf.indexOf(val) >= 0
}

// Sort the list.
func (sf *LinkedList[T]) Sort(less func(a, b T) int) {
	if sf.Len() <= 1 {
		return
	}

	// get all the Values and sort the data
	vs := sf.Values()
	slices.SortFunc(vs, less)

	// clear the linked list and push it back
	sf.Clear()
	for i := 0; i < len(vs); i++ {
		sf.PushBack(vs[i])
	}
}

// Values get a copy of all the values in the list.
func (sf *LinkedList[T]) Values() []T {
	if sf.Len() == 0 {
		return []T{}
	}

	values := make([]T, 0, sf.Len())
	sf.Iterator(func(v T) bool {
		values = append(values, v)
		return true
	})
	return values
}

// getElement returns the element at the specified position.
func (sf *LinkedList[T]) getElement(index int) *list.Element[T] {
	var e *list.Element[T]

	if i, length := 0, sf.Len(); index < (length >> 1) {
		for i, e = 0, sf.l.Front(); i < index; i++ {
			e = e.Next()
		}
	} else {
		for i, e = length-1, sf.l.Back(); i > index; i-- {
			e = e.Prev()
		}
	}
	return e
}

// indexOf returns the index of the first occurrence of the specified element
// in this list, or -1 if this list does not contain the element.
func (sf *LinkedList[T]) indexOf(val T) int {
	for index, e := 0, sf.l.Front(); e != nil; e = e.Next() {
		if val == e.Value {
			return index
		}
		index++
	}
	return -1
}
