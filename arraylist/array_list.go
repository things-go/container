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

	"github.com/things-go/container"
	"github.com/things-go/container/comparator"
)

var _ container.List = (*List)(nil)

// List represents an array list.
// It implements the interface list.Interface.
type List struct {
	items []interface{}
	cmp   comparator.Comparator
}

// Option option for New.
type Option func(l *List)

// WithComparator with user's Comparator.
func WithComparator(cmp comparator.Comparator) Option {
	return func(l *List) {
		l.cmp = cmp
	}
}

// New initializes and returns an ArrayList.
func New(opts ...Option) *List {
	l := &List{
		items: []interface{}{},
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (sf *List) Len() int { return len(sf.items) }

// IsEmpty returns the list l is empty or not.
func (sf *List) IsEmpty() bool { return sf.Len() == 0 }

// Clear initializes or clears list l.
func (sf *List) Clear() { sf.items = make([]interface{}, 0) }

// Push inserts a new element e with value v at the back of list l.
func (sf *List) Push(items interface{}) { sf.items = append(sf.items, items) }

// PushFront inserts a new element e with value v at the front of list l.
func (sf *List) PushFront(v interface{}) {
	sf.items = append(sf.items, v)
	moveLastToFirst(sf.items)
}

// PushBack inserts a new element e with value v at the back of list l.
func (sf *List) PushBack(v interface{}) { sf.items = append(sf.items, v) }

// Add inserts the specified element at the specified position in this list.
func (sf *List) Add(index int, val interface{}) error {
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
func (sf *List) PushFrontList(other *List) {
	items := make([]interface{}, 0, len(sf.items)+len(other.items))
	items = append(items, other.items...)
	items = append(items, sf.items...)
	sf.items = items
}

// PushBackList inserts a copy of an other list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (sf *List) PushBackList(other *List) {
	sf.items = append(sf.items, other.items...)
}

// Poll return the front element value and then remove from list.
func (sf *List) Poll() interface{} {
	return sf.PollFront()
}

// PollFront return the front element value and then remove from list.
func (sf *List) PollFront() interface{} {
	var val interface{}

	if n := len(sf.items); n > 0 {
		moveFirstToLast(sf.items)
		val = sf.items[n-1]
		sf.items[n-1] = nil // for gc
		sf.items = sf.items[:n-1]
	}
	return val
}

// PollBack return the back element value and then remove from list.
func (sf *List) PollBack() interface{} {
	var val interface{}

	if n := len(sf.items); n > 0 {
		val = sf.items[n-1]
		sf.items[n-1] = nil // for gc
		sf.items = sf.items[:n-1]
	}
	return val
}

// Remove removes the element at the specified position in this list.
// It returns an error if the index is out of range.
func (sf *List) Remove(index int) (interface{}, error) {
	if index < 0 || index >= len(sf.items) {
		return nil, fmt.Errorf("index out of range, index:%d, len:%d", index, sf.Len())
	}

	val := sf.items[index]
	// sf.items = append(sf.items[:index], sf.items[(index+1):]...)
	moveLastToFirst(sf.items[index:])
	sf.items[len(sf.items)-1] = nil
	sf.items = sf.items[:len(sf.items)-1]
	sf.shrinkList()
	return val, nil
}

// RemoveValue removes the first occurrence of the specified element from this list, if it is present.
// It returns false if the target value isn't present, otherwise returns true.
func (sf *List) RemoveValue(val interface{}) bool {
	if sf.Len() == 0 {
		return false
	}

	if idx := sf.indexOf(val); idx >= 0 {
		// sf.items = append(sf.items[:idx], sf.items[(idx+1):]...)
		moveLastToFirst(sf.items[idx:])
		sf.items[len(sf.items)-1] = nil
		sf.items = sf.items[:len(sf.items)-1]
		sf.shrinkList()
		return true
	}
	return false
}

// Get returns the element at the specified position in this list. The index must be in the range of [0, size).
func (sf *List) Get(index int) (interface{}, error) {
	if index < 0 || index >= len(sf.items) {
		return nil, fmt.Errorf("index out of range, index:%d, len:%d", index, sf.Len())
	}

	return sf.items[index], nil
}

// Peek return the front element value.
func (sf *List) Peek() interface{} {
	return sf.PeekFront()
}

// PeekFront return the front element value.
func (sf *List) PeekFront() interface{} {
	if len(sf.items) > 0 {
		return sf.items[0]
	}
	return nil
}

// PeekBack return the back element value.
func (sf *List) PeekBack() interface{} {
	if len(sf.items) > 0 {
		return sf.items[len(sf.items)-1]
	}
	return nil
}

// Iterator returns an iterator over the elements in this list in proper sequence.
func (sf *List) Iterator(f func(interface{}) bool) {
	for index := 0; index < sf.Len(); index++ {
		if f == nil || !f(sf.items[index]) {
			return
		}
	}
}

// ReverseIterator returns an iterator over the elements in this list in reverse sequence as Iterator.
func (sf *List) ReverseIterator(f func(interface{}) bool) {
	for index := sf.Len() - 1; index >= 0; index-- {
		if f == nil || !f(sf.items[index]) {
			return
		}
	}
}

// Contains contains the value.
func (sf *List) Contains(val interface{}) bool {
	return val != nil && sf.indexOf(val) >= 0
}

// Sort sort the list.
func (sf *List) Sort(reverse ...bool) {
	if sf.Len() <= 1 {
		return
	}
	ct := comparator.NewContainer(sf.items, sf.cmp)
	if len(reverse) > 0 && reverse[0] {
		ct.Reverse()
	}
	ct.Sort()
}

// Values get a copy of all the values in the list.
func (sf *List) Values() []interface{} {
	items := make([]interface{}, 0, len(sf.items))
	items = append(items, sf.items...)
	return items
}

func (sf *List) shrinkList() {
	oldLen, oldCap := len(sf.items), cap(sf.items)
	if oldCap > 1024 && oldLen <= oldCap/4 { // shrink when len(list) <= cap(list)/4
		newItems := make([]interface{}, oldLen)
		copy(newItems, sf.items)
		sf.Clear()
		sf.items = newItems
	}
}

// indexOf returns the index of the first occurrence of the specified element
// in this list, or -1 if this list does not contain the element.
func (sf *List) indexOf(val interface{}) int {
	for i, v := range sf.items {
		if sf.compare(v, val) {
			return i
		}
	}
	return -1
}

func (sf *List) compare(v1, v2 interface{}) bool {
	if sf.cmp != nil {
		return sf.cmp.Compare(v1, v2) == 0
	}
	return comparator.Compare(v1, v2) == 0
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
