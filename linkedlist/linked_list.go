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

// Package linkedlist implements both an linked and a list.
package linkedlist

import (
	"container/list"
	"fmt"

	"github.com/things-go/container"
	"github.com/things-go/container/comparator"
)

var _ container.List = (*LinkedList)(nil)

// LinkedList represents a doubly linked list.
// It implements the interface list.Interface.
type LinkedList struct {
	l   *list.List
	cmp comparator.Comparator
}

// Option option for New.
type Option func(l *LinkedList)

// WithComparator with user's Comparator.
func WithComparator(cmp comparator.Comparator) Option {
	return func(l *LinkedList) {
		l.cmp = cmp
	}
}

// New initializes and returns an LinkedList.
func New(opts ...Option) *LinkedList {
	l := &LinkedList{l: list.New()}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (sf *LinkedList) Len() int { return sf.l.Len() }

// IsEmpty returns the list l is empty or not.
func (sf *LinkedList) IsEmpty() bool { return sf.l.Len() == 0 }

// Clear initializes or clears list l.
func (sf *LinkedList) Clear() { sf.l.Init() }

// Push inserts a new element e with value v at the back of list l.
func (sf *LinkedList) Push(v interface{}) { sf.l.PushBack(v) }

// PushFront inserts a new element e with value v at the front of list l.
func (sf *LinkedList) PushFront(v interface{}) { sf.l.PushFront(v) }

// PushBack inserts a new element e with value v at the back of list l.
func (sf *LinkedList) PushBack(v interface{}) { sf.l.PushBack(v) }

// Add add to the index of the list with value.
func (sf *LinkedList) Add(index int, val interface{}) error {
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

// PushFrontList inserts a copy of an other list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (sf *LinkedList) PushFrontList(other *LinkedList) {
	sf.l.PushFrontList(other.l)
}

// PushBackList inserts a copy of an other list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (sf *LinkedList) PushBackList(other *LinkedList) {
	sf.l.PushBackList(other.l)
}

// Poll return the front element value and then remove from list.
func (sf *LinkedList) Poll() interface{} {
	return sf.PollFront()
}

// PollFront return the front element value and then remove from list.
func (sf *LinkedList) PollFront() interface{} {
	e := sf.l.Front()
	if e != nil {
		return sf.l.Remove(e)
	}
	return nil
}

// PollBack return the back element value and then remove from list.
func (sf *LinkedList) PollBack() interface{} {
	e := sf.l.Back()
	if e != nil {
		return sf.l.Remove(e)
	}
	return nil
}

// Remove remove the index in the list.
func (sf *LinkedList) Remove(index int) (interface{}, error) {
	if index < 0 || index >= sf.Len() {
		return nil, fmt.Errorf("index out of range, index:%d, len:%d", index, sf.Len())
	}
	return sf.l.Remove(sf.getElement(index)), nil
}

// RemoveValue remove the value in the list.
func (sf *LinkedList) RemoveValue(val interface{}) bool {
	if sf.Len() == 0 {
		return false
	}

	for e := sf.l.Front(); e != nil; e = e.Next() {
		if sf.compare(val, e.Value) {
			sf.l.Remove(e)
			return true
		}
	}
	return false
}

// Get get the index in the list.
func (sf *LinkedList) Get(index int) (interface{}, error) {
	if index < 0 || index >= sf.Len() {
		return nil, fmt.Errorf("index out of range, index: %d, len: %d", index, sf.Len())
	}
	return sf.getElement(index).Value, nil
}

// Peek return the front element value.
func (sf *LinkedList) Peek() interface{} {
	return sf.PeekFront()
}

// PeekFront return the front element value.
func (sf *LinkedList) PeekFront() interface{} {
	if e := sf.l.Front(); e != nil {
		return e.Value
	}
	return nil
}

// PeekBack return the back element value.
func (sf *LinkedList) PeekBack() interface{} {
	if e := sf.l.Back(); e != nil {
		return e.Value
	}
	return nil
}

// Iterator iterator the list.
func (sf *LinkedList) Iterator(cb func(interface{}) bool) {
	for e := sf.l.Front(); e != nil; e = e.Next() {
		if cb == nil || !cb(e.Value) {
			return
		}
	}
}

// ReverseIterator reverse iterator the list.
func (sf *LinkedList) ReverseIterator(cb func(interface{}) bool) {
	for e := sf.l.Back(); e != nil; e = e.Prev() {
		if cb == nil || !cb(e.Value) {
			return
		}
	}
}

// Contains contains the value.
func (sf *LinkedList) Contains(val interface{}) bool {
	return val != nil && sf.indexOf(val) >= 0
}

// Sort sort the list.
func (sf *LinkedList) Sort(reverse ...bool) {
	if sf.Len() <= 1 {
		return
	}

	// get all the Values and sort the data
	vs := sf.Values()
	ct := comparator.NewContainer(vs, sf.cmp)
	if len(reverse) > 0 && reverse[0] {
		ct.Reverse()
	}
	ct.Sort()

	// clear the linked list and push it back
	sf.Clear()
	for i := 0; i < len(vs); i++ {
		sf.PushBack(vs[i])
	}
}

// Values get a copy of all the values in the list.
func (sf *LinkedList) Values() []interface{} {
	if sf.Len() == 0 {
		return []interface{}{}
	}

	values := make([]interface{}, 0, sf.Len())
	sf.Iterator(func(v interface{}) bool {
		values = append(values, v)
		return true
	})
	return values
}

// getElement returns the element at the specified position.
func (sf *LinkedList) getElement(index int) *list.Element {
	var e *list.Element

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
func (sf *LinkedList) indexOf(val interface{}) int {
	for index, e := 0, sf.l.Front(); e != nil; e = e.Next() {
		if sf.compare(val, e.Value) {
			return index
		}
		index++
	}
	return -1
}

func (sf *LinkedList) compare(v1, v2 interface{}) bool {
	if sf.cmp != nil {
		return sf.cmp.Compare(v1, v2) == 0
	}
	return comparator.Compare(v1, v2) == 0
}
