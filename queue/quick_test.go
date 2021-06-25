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

package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuickQueueLen(t *testing.T) {
	q := NewQuickQueue()
	q.Add(5)
	q.Add(6)
	assert.Equal(t, 2, q.Len())
}

func TestQuickQueuePeek(t *testing.T) {
	q := NewQuickQueue()
	q.Add(5)
	q.Add("hello")

	val1, ok := q.Peek().(int)
	assert.True(t, ok)
	assert.Equal(t, 5, val1)

	val2, ok := q.Peek().(int)
	assert.True(t, ok)
	assert.Equal(t, 5, val2)

	q.Poll()
	q.Poll()

	val3 := q.Peek()
	assert.Nil(t, val3)
}

func TestQuickQueueValue(t *testing.T) {
	q := NewQuickQueue()
	q.Add(15)
	q.Add(11)
	q.Add(19)
	q.Add(12)
	q.Add(8)
	q.Add(1)

	// len
	require.Equal(t, 6, q.Len())

	// Peek
	require.Equal(t, 15, q.Peek())
	require.Equal(t, 6, q.Len())

	// Contains
	require.True(t, q.Contains(19))
	require.False(t, q.Contains(10000))

	// Poll
	require.Equal(t, 15, q.Poll())
	require.Nil(t, q.head[0])
	require.Equal(t, 5, q.Len())

	// add new to tail
	q.Add(23)
	q.Add(3)
	q.Add(17)
	q.Add(7)

	require.Equal(t, 9, q.Len())

	// Contains (again)
	require.True(t, q.Contains(19))
	require.False(t, q.Contains(10000))

	// Remove
	q.Remove(19)
	require.False(t, q.Contains(19))
	require.Nil(t, q.head[1])
	require.Equal(t, 4, len(q.head)-q.headPos)

	// Remove
	q.Remove(8)
	require.False(t, q.Contains(8))
	require.Equal(t, 3, len(q.head)-q.headPos)

	q.Remove(17)
	require.False(t, q.Contains(17))
	require.Equal(t, 3, len(q.tail))

	require.Equal(t, 6, q.Len())
}

func TestQuickQueuePoll(t *testing.T) {
	q := NewQuickQueue()

	q.Add(5)
	q.Add("hello")
	val1, ok := q.Poll().(int)
	assert.True(t, ok)
	assert.Equal(t, 5, val1)

	val2, ok := q.Poll().(string)
	assert.True(t, ok)
	assert.Equal(t, "hello", val2)

	val3 := q.Poll()
	assert.Nil(t, val3)
}

func TestQuickQueueIsEmpty(t *testing.T) {
	q := NewQuickQueue()
	q.Add(5)
	q.Add(6)
	assert.False(t, q.IsEmpty())

	assert.Equal(t, 5, q.Poll())
	assert.Equal(t, 6, q.Peek())
	q.Clear()
	assert.Equal(t, 0, q.Len())
	assert.True(t, q.IsEmpty())
}

func TestQuickQueueInit(t *testing.T) {
	q := NewQuickQueue()
	q.Add(5)
	q.Add(6)
	q.Clear()

	assert.Equal(t, 0, q.Len())
}

func TestQuickQueueWithComparator(t *testing.T) {
	q := NewQuickQueue(WithComparator(&student{}))

	q.Add(&student{name: "benjamin", age: 34})
	q.Add(&student{name: "alice", age: 21})
	q.Add(&student{name: "john", age: 42})
	q.Add(&student{name: "roy", age: 28})
	q.Add(&student{name: "moss", age: 25})

	assert.Equal(t, 5, q.Len())

	assert.True(t, q.Contains(&student{name: "alice", age: 21}))

	// Peek
	v, ok := q.Peek().(*student)
	require.True(t, ok)
	require.True(t, v.name == "benjamin" && v.age == 34)

	v, ok = q.Poll().(*student)
	require.True(t, ok)
	require.True(t, v.name == "benjamin" && v.age == 34)

	v, ok = q.Poll().(*student)
	require.True(t, ok)
	require.True(t, v.name == "alice" && v.age == 21)

	v, ok = q.Poll().(*student)
	require.True(t, ok)
	require.True(t, v.name == "john" && v.age == 42)

	v, ok = q.Poll().(*student)
	require.True(t, ok)
	require.True(t, v.name == "roy" && v.age == 28)

	v, ok = q.Poll().(*student)
	require.True(t, ok)
	require.True(t, v.name == "moss" && v.age == 25)

	// The queue should be empty now
	require.Zero(t, q.Len())
	require.Nil(t, q.Peek())
	require.Nil(t, q.Poll())
}
