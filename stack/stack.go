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

// Package stack implements a Stack, which orders elements in a LIFO (last-in-first-out) manner.
package stack

import (
	"github.com/things-go/container"
	"github.com/things-go/container/go/list"
)

var _ container.Stack[int] = (*Stack[int])(nil)

// Stack is LIFO implement list.List.
type Stack[T any] struct {
	list *list.List[T]
}

// New creates a Stack. which implement interface stack.Interface.
func New[T any]() *Stack[T] { return &Stack[T]{list.New[T]()} }

// Len returns the length of this priority queue.
func (qs *Stack[T]) Len() int { return qs.list.Len() }

// IsEmpty returns true if this Stack contains no elements.
func (qs *Stack[T]) IsEmpty() bool { return qs.list.Len() == 0 }

// Clear removes all the elements from this Stack.
func (qs *Stack[T]) Clear() { qs.list.Init() }

// Push pushes an element into this Stack.
func (qs *Stack[T]) Push(val T) { qs.list.PushFront(val) }

// Pop pops the element on the top of this Stack.
func (qs *Stack[T]) Pop() (v T, ok bool) {
	if e := qs.list.Front(); e != nil {
		return qs.list.Remove(e), true
	}
	return v, false
}

// Peek retrieves, but does not remove,
// the element on the top of this Stack, or return nil if this Stack is empty.
func (qs *Stack[T]) Peek() (v T, ok bool) {
	if e := qs.list.Front(); e != nil {
		return e.Value, true
	}
	return v, false
}

// Copy returns a copy of this stack.
func (qs *Stack[T]) Clone() *Stack[T] {
	s := New[T]()
	for e := qs.list.Front(); e != nil; e = e.Next() {
		s.Push(e.Value)
	}
	return s
}
