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

package priorityqueue

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPriorityQueueLen(t *testing.T) {
	// init 3 elements
	q := New(WithItems([]interface{}{5, 6, 7}))

	require.Equal(t, 3, q.Len())
	require.False(t, q.IsEmpty())

	// remove one element
	q.Remove(6)
	require.Equal(t, 2, q.Len())

	// remove one element not exist
	q.Remove(10000)
	require.Equal(t, 2, q.Len())

	// Clear all elements
	q.Clear()
	require.Zero(t, q.Len())
	require.True(t, q.IsEmpty())

	// remove one element if not any element in queue
	q.Remove(10000)
}

func TestPriorityQueueValue(t *testing.T) {
	// create priority queue
	q := New()
	q.Add(15)
	q.Add(19)
	q.Add(12)
	q.Add(8)
	q.Add(13)

	require.Equal(t, 5, q.Len())

	// Peek
	require.Equal(t, 8, q.Peek())
	require.Equal(t, 5, q.Len())

	// Contains
	require.True(t, q.Contains(12))
	require.False(t, q.Contains(10000))

	// Poll
	require.Equal(t, 8, q.Poll())
	require.Equal(t, 4, q.Len())

	require.Equal(t, 12, q.Poll())
	require.Equal(t, 3, q.Len())

	// Contains (again)
	require.False(t, q.Contains(12))
	require.False(t, q.Contains(10000))

	// Remove
	require.True(t, q.Contains(15))
	q.Remove(15)
	require.False(t, q.Contains(15))
}

func TestPriorityQueueMinHeap(t *testing.T) {
	pq := New()
	pqTestPriorityQueueSortImpl(t, pq, []interface{}{15, 19, 12, 8, 13}, []interface{}{8, 12, 13, 15, 19})
}

func TestPriorityQueueMinHeapWithComparator(t *testing.T) {
	pq := New(WithComparator(CompareMyInt))
	pqTestPriorityQueueSortImpl(t, pq, []interface{}{15, 19, 12, 8, 13}, []interface{}{19, 15, 13, 12, 8})
}

func TestPriorityQueueMaxHeap(t *testing.T) {
	pq := New(WithMaxHeap())
	pqTestPriorityQueueSortImpl(t, pq, []interface{}{15, 19, 12, 8, 13}, []interface{}{19, 15, 13, 12, 8})
}

func TestPriorityQueueMaxHeapWithComparator(t *testing.T) {
	q := New(WithComparator(CompareMyInt), WithMaxHeap())
	pqTestPriorityQueueSortImpl(t, q, []interface{}{15, 19, 12, 8, 13}, []interface{}{8, 12, 13, 15, 19})
}

func pqTestPriorityQueueSortImpl(t *testing.T, q *Queue, input, expected []interface{}) {
	for i := 0; i < len(input); i++ {
		q.Add(input[i])
	}

	require.Equal(t, len(input), q.Len())
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], q.Poll())
	}
	require.Zero(t, q.Len())
}

func TestPriorityQueueDeleteMinHeap(t *testing.T) {
	pq := New()
	pqTestPriorityQueueDeleteImpl(t, pq, []interface{}{15, 19, 12, 8, 13}, []interface{}{8, 12, 13, 15}, 19)
}

func TestPriorityQueueDeleteMinHeapWithComparator(t *testing.T) {
	pq := New(WithComparator(CompareMyInt))
	pqTestPriorityQueueDeleteImpl(t, pq, []interface{}{15, 19, 12, 8, 13}, []interface{}{19, 13, 12, 8}, 15)
}

func TestPriorityQueueDeleteMaxHeap(t *testing.T) {
	pq := New(WithMaxHeap())
	pqTestPriorityQueueDeleteImpl(t, pq, []interface{}{15, 19, 12, 8, 13}, []interface{}{19, 15, 13, 8}, 12)
}

func TestPriorityQueueDeleteMaxHeapWithComparator(t *testing.T) {
	pq := New(WithComparator(CompareMyInt), WithMaxHeap())
	pqTestPriorityQueueDeleteImpl(t, pq, []interface{}{15, 19, 12, 8, 13}, []interface{}{12, 13, 15, 19}, 8)
}

func pqTestPriorityQueueDeleteImpl(t *testing.T, q *Queue, input, expected []interface{}, val interface{}) {
	for i := 0; i < len(input); i++ {
		q.Add(input[i])
	}

	q.Remove(val)
	require.Equal(t, len(input)-1, q.Len())
	assert.False(t, q.Contains(val))
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], q.Poll())
	}
	require.Zero(t, q.Len())
}

// Compare returns reverse order.
func CompareMyInt(v1, v2 interface{}) int {
	i1, i2 := v1.(int), v2.(int)
	if i1 < i2 {
		return 1
	}
	if i1 > i2 {
		return -1
	}
	return 0
}

func TestPQComparator(t *testing.T) {
	pq := New(WithComparator(CompareStudent))

	pq.Add(&student{name: "benjamin", age: 34})
	pq.Add(&student{name: "alice", age: 21})
	pq.Add(&student{name: "john", age: 42})
	pq.Add(&student{name: "roy", age: 28})
	pq.Add(&student{name: "moss", age: 25})

	assert.Equal(t, 5, pq.Len())

	// Peek
	v, ok := pq.Peek().(*student)
	require.True(t, ok)
	require.True(t, v.name == "john" && v.age == 42)

	v, ok = pq.Poll().(*student)
	require.True(t, ok)
	require.True(t, v.name == "john" && v.age == 42)

	v, ok = pq.Poll().(*student)
	require.True(t, ok)
	require.True(t, v.name == "benjamin" && v.age == 34)

	v, ok = pq.Poll().(*student)
	require.True(t, ok)
	require.True(t, v.name == "roy" && v.age == 28)

	v, ok = pq.Poll().(*student)
	require.True(t, ok)
	require.True(t, v.name == "moss" && v.age == 25)

	v, ok = pq.Poll().(*student)
	require.True(t, ok)
	require.True(t, v.name == "alice" && v.age == 21)

	// The queue should be empty now
	require.Zero(t, pq.Len())
	require.Nil(t, pq.Peek())
	require.Nil(t, pq.Poll())
}

type student struct {
	name string
	age  int
}

// Compare returns -1, 0 or 1 when the first student's age is greater, equal to, or less than the second student's age.
func CompareStudent(v1, v2 interface{}) int {
	s1, s2 := v1.(*student), v2.(*student)
	if s1.age < s2.age {
		return 1
	}
	if s1.age > s2.age {
		return -1
	}
	return 0
}
