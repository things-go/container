// Copyright [2020] [thinkgos]
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

	"golang.org/x/exp/slices"

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
func (sf *List[T]) Len() int { return len(sf.items) }

// IsEmpty returns the list l is empty or not.
func (sf *List[T]) IsEmpty() bool { return sf.Len() == 0 }

// Clear initializes or clears list l.
func (sf *List[T]) Clear() { sf.items = make([]T, 0) }

// Push inserts a new element e with value v at the back of list l.
func (sf *List[T]) Push(items T) { sf.items = append(sf.items, items) }

// PushFront inserts a new element e with value v at the front of list l.
func (sf *List[T]) PushFront(v T) {
	sf.items = append(sf.items, v)
	moveLastToFirst(sf.items)
}

// PushBack inserts a new element e with value v at the back of list l.
func (sf *List[T]) PushBack(v T) { sf.items = append(sf.items, v) }

// Add inserts the specified element at the specified position in this list.
func (sf *List[T]) Add(index int, val T) error {
	if index < 0 || index > len(sf.items) {
		return fmt.Errorf("index out of range, index:%d, len:%d", index, sf.Len())
	}

	if index == sf.Len() {
		sf.Push(val)
	} else {
		length := len(sf.items)
		sf.items = append(sf.items, val)
		copy(sf.items[index+1:], sf.items[index:length])
		sf.items[index] = val
	}
	return nil
}

// PushFrontList inserts a copy of an other list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (sf *List[T]) PushFrontList(other *List[T]) {
	items := make([]T, 0, len(sf.items)+len(other.items))
	items = append(items, other.items...)
	items = append(items, sf.items...)
	sf.items = items
}

// PushBackList inserts a copy of an other list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (sf *List[T]) PushBackList(other *List[T]) {
	sf.items = append(sf.items, other.items...)
}

// Poll return the front element value and then remove from list.
func (sf *List[T]) Poll() (val T, ok bool) {
	return sf.PollFront()
}

// PollFront return the front element value and then remove from list.
func (sf *List[T]) PollFront() (val T, ok bool) {
	var placeholder T

	if n := len(sf.items); n > 0 {
		moveFirstToLast(sf.items)
		val = sf.items[n-1]
		sf.items[n-1] = placeholder // for gc
		sf.items = sf.items[:n-1]
		ok = true
	}
	return val, ok
}

// PollBack return the back element value and then remove from list.
func (sf *List[T]) PollBack() (val T, ok bool) {
	var placeholder T

	if n := len(sf.items); n > 0 {
		val = sf.items[n-1]
		sf.items[n-1] = placeholder // for gc
		sf.items = sf.items[:n-1]
		ok = true
	}
	return val, ok
}

// Remove removes the element at the specified position in this list.
// It returns an error if the index is out of range.
func (sf *List[T]) Remove(index int) (val T, err error) {
	var placeholder T

	if index < 0 || index >= len(sf.items) {
		return val, fmt.Errorf("index out of range, index:%d, len:%d", index, sf.Len())
	}

	val = sf.items[index]
	// sf.items = append(sf.items[:index], sf.items[(index+1):]...)
	moveLastToFirst(sf.items[index:])
	sf.items[len(sf.items)-1] = placeholder
	sf.items = sf.items[:len(sf.items)-1]
	sf.shrinkList()
	return val, nil
}

// RemoveValue removes the first occurrence of the specified element from this list, if it is present.
// It returns false if the target value isn't present, otherwise returns true.
func (sf *List[T]) RemoveValue(val T) bool {
	var placeholder T

	if sf.Len() == 0 {
		return false
	}

	if idx := sf.indexOf(val); idx >= 0 {
		// sf.items = append(sf.items[:idx], sf.items[(idx+1):]...)
		moveLastToFirst(sf.items[idx:])
		sf.items[len(sf.items)-1] = placeholder
		sf.items = sf.items[:len(sf.items)-1]
		sf.shrinkList()
		return true
	}
	return false
}

// Get returns the element at the specified position in this list. The index must be in the range of [0, size).
func (sf *List[T]) Get(index int) (val T, err error) {
	if index < 0 || index >= len(sf.items) {
		return val, fmt.Errorf("index out of range, index:%d, len:%d", index, sf.Len())
	}

	return sf.items[index], nil
}

// Peek return the front element value.
func (sf *List[T]) Peek() (val T, ok bool) {
	return sf.PeekFront()
}

// PeekFront return the front element value.
func (sf *List[T]) PeekFront() (val T, ok bool) {
	if len(sf.items) > 0 {
		return sf.items[0], true
	}
	return val, false
}

// PeekBack return the back element value.
func (sf *List[T]) PeekBack() (val T, ok bool) {
	if len(sf.items) > 0 {
		return sf.items[len(sf.items)-1], true
	}
	return val, false
}

// Iterator returns an iterator over the elements in this list in proper sequence.
func (sf *List[T]) Iterator(f func(T) bool) {
	for index := 0; index < sf.Len(); index++ {
		if f == nil || !f(sf.items[index]) {
			return
		}
	}
}

// ReverseIterator returns an iterator over the elements in this list in reverse sequence as Iterator.
func (sf *List[T]) ReverseIterator(f func(T) bool) {
	for index := sf.Len() - 1; index >= 0; index-- {
		if f == nil || !f(sf.items[index]) {
			return
		}
	}
}

// Contains contains the value.
func (sf *List[T]) Contains(val T) bool {
	return sf.indexOf(val) >= 0
}

// Sort the list.
func (sf *List[T]) Sort(less func(a, b T) bool) {
	if sf.Len() <= 1 {
		return
	}
	slices.SortFunc(sf.items, less)
}

// Values get a copy of all the values in the list.
func (sf *List[T]) Values() []T {
	items := make([]T, 0, len(sf.items))
	items = append(items, sf.items...)
	return items
}

func (sf *List[T]) shrinkList() {
	oldLen, oldCap := len(sf.items), cap(sf.items)
	if oldCap > 1024 && oldLen <= oldCap/4 { // shrink when len(list) <= cap(list)/4
		newItems := make([]T, oldLen)
		copy(newItems, sf.items)
		sf.Clear()
		sf.items = newItems
	}
}

// indexOf returns the index of the first occurrence of the specified element
// in this list, or -1 if this list does not contain the element.
func (sf *List[T]) indexOf(val T) int {
	for i, v := range sf.items {
		if v == val {
			return i
		}
	}
	return -1
}

func moveLastToFirst[T any](items []T) {
	for i := 0; i < len(items); i++ {
		items[i], items[len(items)-1] = items[len(items)-1], items[i]
	}
}

func moveFirstToLast[T any](items []T) {
	for i := 0; i < len(items); i++ {
		items[0], items[len(items)-1-i] = items[len(items)-1-i], items[0]
	}
}
