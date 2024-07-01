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
	"github.com/things-go/container/go/list"
)

var _ container.List[int] = (*LinkedList[int])(nil)

// LinkedList represents a doubly linked list.
// It implements the interface list.Interface.
type LinkedList[T comparable] struct {
	list *list.List[T]
}

// New initializes and returns an LinkedList.
func New[T comparable]() *LinkedList[T] {
	return &LinkedList[T]{list: list.New[T]()}
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (ll *LinkedList[T]) Len() int { return ll.list.Len() }

// IsEmpty returns the list l is empty or not.
func (ll *LinkedList[T]) IsEmpty() bool { return ll.list.Len() == 0 }

// Clear initializes or clears list l.
func (ll *LinkedList[T]) Clear() { ll.list.Init() }

// Push inserts a new element e with value v at the back of list l.
func (ll *LinkedList[T]) Push(v T) { ll.list.PushBack(v) }

// PushFront inserts a new element e with value v at the front of list l.
func (ll *LinkedList[T]) PushFront(v T) { ll.list.PushFront(v) }

// PushBack inserts a new element e with value v at the back of list l.
func (ll *LinkedList[T]) PushBack(v T) { ll.list.PushBack(v) }

// Add add to the index of the list with value.
func (ll *LinkedList[T]) Add(index int, val T) error {
	if index < 0 || index > ll.Len() {
		return fmt.Errorf("index out of range, index: %d, len: %d", index, ll.Len())
	}

	if index == ll.Len() {
		ll.list.PushBack(val)
	} else {
		ll.list.InsertBefore(val, ll.getElement(index))
	}
	return nil
}

// PushFrontList inserts a copy of another list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (ll *LinkedList[T]) PushFrontList(other *LinkedList[T]) {
	ll.list.PushFrontList(other.list)
}

// PushBackList inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (ll *LinkedList[T]) PushBackList(other *LinkedList[T]) {
	ll.list.PushBackList(other.list)
}

// Poll return the front element value and then remove from list.
func (ll *LinkedList[T]) Poll() (T, bool) {
	return ll.PollFront()
}

// PollFront return the front element value and then remove from list.
func (ll *LinkedList[T]) PollFront() (val T, ok bool) {
	e := ll.list.Front()
	if e != nil {
		return ll.list.Remove(e), true
	}
	return val, false
}

// PollBack return the back element value and then remove from list.
func (ll *LinkedList[T]) PollBack() (val T, ok bool) {
	e := ll.list.Back()
	if e != nil {
		return ll.list.Remove(e), true
	}
	return val, false
}

// Remove remove the index in the list.
func (ll *LinkedList[T]) Remove(index int) (val T, err error) {
	if index < 0 || index >= ll.Len() {
		return val, fmt.Errorf("index out of range, index:%d, len:%d", index, ll.Len())
	}
	return ll.list.Remove(ll.getElement(index)), nil
}

// RemoveValue remove the value in the list.
func (ll *LinkedList[T]) RemoveValue(val T) bool {
	if ll.Len() == 0 {
		return false
	}

	for e := ll.list.Front(); e != nil; e = e.Next() {
		if val == e.Value {
			ll.list.Remove(e)
			return true
		}
	}
	return false
}

// Get the index in the list.
func (ll *LinkedList[T]) Get(index int) (val T, err error) {
	if index < 0 || index >= ll.Len() {
		return val, fmt.Errorf("index out of range, index: %d, len: %d", index, ll.Len())
	}
	return ll.getElement(index).Value, nil
}

// Peek return the front element value.
func (ll *LinkedList[T]) Peek() (T, bool) {
	return ll.PeekFront()
}

// PeekFront return the front element value.
func (ll *LinkedList[T]) PeekFront() (val T, ok bool) {
	if e := ll.list.Front(); e != nil {
		return e.Value, true
	}
	return val, false
}

// PeekBack return the back element value.
func (ll *LinkedList[T]) PeekBack() (val T, ok bool) {
	if e := ll.list.Back(); e != nil {
		return e.Value, true
	}
	return val, false
}

// Iterator the list.
func (ll *LinkedList[T]) Iterator(cb func(T) bool) {
	for e := ll.list.Front(); e != nil; e = e.Next() {
		if cb == nil || !cb(e.Value) {
			return
		}
	}
}

// ReverseIterator reverse iterator the list.
func (ll *LinkedList[T]) ReverseIterator(cb func(T) bool) {
	for e := ll.list.Back(); e != nil; e = e.Prev() {
		if cb == nil || !cb(e.Value) {
			return
		}
	}
}

// Contains contains the value.
func (ll *LinkedList[T]) Contains(val T) bool {
	return ll.indexOf(val) >= 0
}

// Sort the list.
func (ll *LinkedList[T]) Sort(less func(a, b T) int) {
	if ll.Len() <= 1 {
		return
	}

	// get all the Values and sort the data
	vs := ll.Values()
	slices.SortFunc(vs, less)

	// clear the linked list and push it back
	ll.Clear()
	for i := 0; i < len(vs); i++ {
		ll.PushBack(vs[i])
	}
}

// Values get a copy of all the values in the list.
func (ll *LinkedList[T]) Values() []T {
	if ll.Len() == 0 {
		return []T{}
	}

	values := make([]T, 0, ll.Len())
	ll.Iterator(func(v T) bool {
		values = append(values, v)
		return true
	})
	return values
}

// getElement returns the element at the specified position.
func (ll *LinkedList[T]) getElement(index int) *list.Element[T] {
	var e *list.Element[T]

	if i, length := 0, ll.Len(); index < (length >> 1) {
		for i, e = 0, ll.list.Front(); i < index; i++ {
			e = e.Next()
		}
	} else {
		for i, e = length-1, ll.list.Back(); i > index; i-- {
			e = e.Prev()
		}
	}
	return e
}

// indexOf returns the index of the first occurrence of the specified element
// in this list, or -1 if this list does not contain the element.
func (ll *LinkedList[T]) indexOf(val T) int {
	for index, e := 0, ll.list.Front(); e != nil; e = e.Next() {
		if val == e.Value {
			return index
		}
		index++
	}
	return -1
}
