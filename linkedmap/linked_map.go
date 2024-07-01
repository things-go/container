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

package linkedmap

import (
	"github.com/things-go/container"
	"github.com/things-go/container/go/list"
)

var _ container.LinkedMap[int, int] = (*LinkedMap[int, int])(nil)

type store[K comparable, V any] struct {
	key   K
	value V
}

// LinkedMap implements the Interface.
type LinkedMap[K comparable, V any] struct {
	data     map[K]*list.Element[*store[K, V]]
	list     *list.List[*store[K, V]]
	capacity int
}

// Option for New.
type Option[K comparable, V any] func(lm *LinkedMap[K, V])

// WithCap with limit capacity.
func WithCap[K comparable, V any](capacity int) Option[K, V] {
	return func(lm *LinkedMap[K, V]) {
		lm.capacity = capacity
	}
}

// New creates a LinkedMap.
func New[K comparable, V any](opts ...Option[K, V]) *LinkedMap[K, V] {
	lm := &LinkedMap[K, V]{
		data: make(map[K]*list.Element[*store[K, V]]),
		list: list.New[*store[K, V]](),
	}
	for _, opt := range opts {
		opt(lm)
	}
	return lm
}

// Cap returns the capacity of elements of list ll.
// The complexity is O(1).
func (lm *LinkedMap[K, V]) Cap() int { return lm.capacity }

// Len returns the number of elements of list ll.
// The complexity is O(1).
func (lm *LinkedMap[K, V]) Len() int { return lm.list.Len() }

// IsEmpty returns the list ll is empty or not.
func (lm *LinkedMap[K, V]) IsEmpty() bool { return lm.Len() == 0 }

// Clear initializes or clears list ll.
func (lm *LinkedMap[K, V]) Clear() {
	lm.data = make(map[K]*list.Element[*store[K, V]])
	lm.list.Init()
}

// Push associates the specified value with the specified key in this map.
// If the map previously contained a mapping for the key,
// the old value is replaced by the specified value. and then move the item to the back of the list.
// If over the capacity, it will remove the back item then push new item to back
// It returns the previous value associated with the specified key, or nil if there was no mapping for the key.
// A nil return can also indicate that the map previously associated nil with the specified key.
func (lm *LinkedMap[K, V]) Push(k K, v V) (V, bool) { return lm.PushBack(k, v) }

// PushFront associates the specified value with the specified key in this map.
// If the map previously contained a mapping for the key,
// the old value is replaced by the specified value. and then move the item to the front of the list.
// If over the capacity, it will remove the back item then push new item to front
// It returns the previous value associated with the specified key, or nil if there was no mapping for the key.
// A nil return can also indicate that the map previously associated nil with the specified key.
func (lm *LinkedMap[K, V]) PushFront(k K, v V) (val V, exist bool) {
	if old, ok := lm.data[k]; ok {
		val = old.Value.value
		old.Value = &store[K, V]{k, v}
		lm.list.MoveToFront(old)
		exist = true
	} else {
		if lm.capacity != 0 && lm.list.Len() >= lm.capacity {
			e := lm.list.Back()
			delete(lm.data, e.Value.key)
			lm.list.Remove(e)
		}
		lm.data[k] = lm.list.PushFront(&store[K, V]{k, v})
	}
	return val, exist
}

// PushBack associates the specified value with the specified key in this map.
// If the map previously contained a mapping for the key,
// the old value is replaced by the specified value. and then move the item to the back of the list.
// If over the capacity, it will remove the back item then push new item to back.
func (lm *LinkedMap[K, V]) PushBack(k K, v V) (val V, exist bool) {
	if old, ok := lm.data[k]; ok {
		val = old.Value.value
		old.Value = &store[K, V]{k, v}
		lm.list.MoveToBack(old)
		exist = true
	} else {
		if lm.capacity != 0 && lm.list.Len() >= lm.capacity {
			e := lm.list.Front()
			delete(lm.data, e.Value.key)
			lm.list.Remove(e)
		}
		lm.data[k] = lm.list.PushBack(&store[K, V]{k, v})
	}
	return val, exist
}

// Poll return the front element value and then remove from list.
func (lm *LinkedMap[K, V]) Poll() (k K, v V, exist bool) { return lm.PollFront() }

// PollFront return the front element value and then remove from list.
func (lm *LinkedMap[K, V]) PollFront() (k K, v V, exist bool) {
	if e := lm.list.Front(); e != nil {
		st := e.Value
		delete(lm.data, st.key)
		lm.list.Remove(e)
		return st.key, st.value, true
	}
	return k, v, false
}

// PollBack return the back element value and then remove from list.
func (lm *LinkedMap[K, V]) PollBack() (k K, v V, exist bool) {
	if e := lm.list.Back(); e != nil {
		st := e.Value
		delete(lm.data, st.key)
		lm.list.Remove(e)
		return st.key, st.value, true
	}
	return k, v, false
}

// Remove removes the mapping for a key from this map if it is present.
// It returns the value to which this map previously associated the key, and true,
// or nil and false if the map contained no mapping for the key.
func (lm *LinkedMap[K, V]) Remove(k K) (val V, exist bool) {
	if oldElement, ok := lm.data[k]; ok {
		val = oldElement.Value.value
		delete(lm.data, k)
		lm.list.Remove(oldElement)
		exist = true
	}
	return val, exist
}

// Contains returns true if this map contains a mapping for the specified key.
func (lm *LinkedMap[K, V]) Contains(k K) bool {
	_, ok := lm.data[k]
	return ok
}

// ContainsValue returns true if this map maps one or more keys to the specified value.
func (lm *LinkedMap[K, V]) ContainsValue(v V, equal func(a, b V) bool) bool {
	for e := lm.list.Front(); e != nil; e = e.Next() {
		if equal(e.Value.value, v) {
			return true
		}
	}
	return false
}

// Get returns the value to which the specified key is mapped, or nil if this map contains no mapping for the key.
func (lm *LinkedMap[K, V]) Get(k K, defaultValue ...V) (val V) {
	if old, ok := lm.data[k]; ok {
		lm.list.MoveToBack(old)
		return old.Value.value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return val
}

// Peek return the front element value .
func (lm *LinkedMap[K, V]) Peek() (k K, v V, exist bool) {
	return lm.PeekFront()
}

// PeekFront return the front element value.
func (lm *LinkedMap[K, V]) PeekFront() (k K, v V, exist bool) {
	if e := lm.list.Front(); e != nil {
		k = e.Value.key
		v = e.Value.value
		exist = true
	}
	return k, v, exist
}

// PeekBack return the back element value .
func (lm *LinkedMap[K, V]) PeekBack() (k K, v V, exist bool) {
	if e := lm.list.Back(); e != nil {
		k = e.Value.key
		v = e.Value.value
		exist = true
	}
	return k, v, exist
}

// Iterator the list.
func (lm *LinkedMap[K, V]) Iterator(cb func(k K, v V) bool) {
	for e := lm.list.Front(); e != nil; e = e.Next() {
		st := e.Value
		if cb == nil || !cb(st.key, st.value) {
			return
		}
	}
}

// ReverseIterator reverse iterator the list.
func (lm *LinkedMap[K, V]) ReverseIterator(cb func(k K, v V) bool) {
	for e := lm.list.Back(); e != nil; e = e.Prev() {
		st := e.Value
		if cb == nil || !cb(st.key, st.value) {
			return
		}
	}
}
