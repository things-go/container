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

func TestQueueLen(t *testing.T) {
	q := New()
	q.Add(5)
	q.Add(6)
	assert.Equal(t, 2, q.Len())
}

func TestQueuePeek(t *testing.T) {
	q := New()
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

func TestQueueValue(t *testing.T) {
	q := New()
	q.Add(15)
	q.Add(19)
	q.Add(12)
	q.Add(8)
	q.Add(1)
	q.Add(13)

	// len
	require.Equal(t, 6, q.Len())

	// Peek
	require.Equal(t, 15, q.Peek())
	require.Equal(t, 6, q.Len())

	// Contains
	require.True(t, q.Contains(12))
	require.False(t, q.Contains(10000))

	// Poll
	require.Equal(t, 15, q.Poll())
	require.Equal(t, 5, q.Len())

	require.Equal(t, 19, q.Poll())
	require.Equal(t, 4, q.Len())

	// Contains (again)
	require.True(t, q.Contains(12))
	require.False(t, q.Contains(10000))

	// Remove
	q.Remove(12)
	require.False(t, q.Contains(12))

	q.Remove(1)
	assert.Nil(t, q.tail.next)
	require.Equal(t, 2, q.Len())
	q.Remove(13)
	assert.Nil(t, q.tail.next)
	require.Equal(t, 1, q.Len())
	q.Remove(8)
	require.Equal(t, 0, q.Len())
	assert.Nil(t, q.head)
	assert.Nil(t, q.tail)
}

func TestQueuePoll(t *testing.T) {
	q := New()

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

func TestQueueIsEmpty(t *testing.T) {
	q := New()
	q.Add(5)
	q.Add(6)
	assert.False(t, q.IsEmpty())

	q.Clear()
	assert.Equal(t, 0, q.Len())
	assert.True(t, q.IsEmpty())
}

func TestQueueInit(t *testing.T) {
	q := New()
	q.Add(5)
	q.Add(6)
	q.Clear()

	assert.Equal(t, 0, q.Len())
}

func TestWithComparator(t *testing.T) {
	q := New(WithComparator(&student{}))

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

type student struct {
	name string
	age  int
}

// Compare returns -1, 0 or 1 when the first student's age is greater, equal to, or less than the second student's age.
func (s *student) Compare(v1, v2 interface{}) int {
	s1, s2 := v1.(*student), v2.(*student)
	if s1.age < s2.age {
		return 1
	}
	if s1.age > s2.age {
		return -1
	}
	return 0
}
