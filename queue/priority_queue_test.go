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

package queue

import (
	"cmp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PriorityQueueLen(t *testing.T) {
	// init 3 elements
	q := NewPriorityQueueWith(false, cmp.Compare, 5, 6, 7)

	require.Equal(t, 3, q.Len())
	require.False(t, q.IsEmpty())

	// remove one element
	t.Log(q.indexOf(6))
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

func Test_PriorityQueueValue(t *testing.T) {
	// create priority queue
	q := NewPriorityQueue[int](false)
	q.Add(15)
	q.Add(19)
	q.Add(12)
	q.Add(8)
	q.Add(13)

	require.Equal(t, 5, q.Len())

	// Peek
	val, ok := q.Peek()
	require.True(t, ok)
	require.Equal(t, int(8), val)
	require.Equal(t, 5, q.Len())

	// Contains
	require.True(t, q.Contains(12))
	require.False(t, q.Contains(10000))

	// Poll
	val, ok = q.Poll()
	require.True(t, ok)
	require.Equal(t, int(8), val)
	require.Equal(t, 4, q.Len())

	val, ok = q.Poll()
	require.True(t, ok)
	require.Equal(t, int(12), val)
	require.Equal(t, 3, q.Len())

	// Contains (again)
	require.False(t, q.Contains(12))
	require.False(t, q.Contains(10000))

	// Remove
	require.True(t, q.Contains(15))
	q.Remove(15)
	require.False(t, q.Contains(15))
}

func Test_PriorityQueueMinHeap(t *testing.T) {
	pq := NewPriorityQueue[int](false)
	pqTestPriorityQueueSortImpl(t, pq, []int{15, 19, 12, 8, 13}, []int{8, 12, 13, 15, 19})
}

func Test_PriorityQueueMaxHeap(t *testing.T) {
	pq := NewPriorityQueue[int](true)
	pqTestPriorityQueueSortImpl(t, pq, []int{15, 19, 12, 8, 13}, []int{19, 15, 13, 12, 8})
}

func pqTestPriorityQueueSortImpl[T comparable](t *testing.T, q *PriorityQueue[T], input, expected []T) {
	for i := 0; i < len(input); i++ {
		q.Add(input[i])
	}

	require.Equal(t, len(input), q.Len())
	for i := 0; i < len(expected); i++ {
		val, ok := q.Poll()
		assert.True(t, ok)
		assert.Equal(t, expected[i], val)
	}
	require.Zero(t, q.Len())
}

func Test_PriorityQueueDeleteMinHeap(t *testing.T) {
	pq := NewPriorityQueue[int](false)
	pqTestPriorityQueueDeleteImpl(t, pq, []int{15, 19, 12, 8, 13}, []int{8, 12, 13, 15}, 19)
}

func Test_PriorityQueueDeleteMinHeapWithComparator(t *testing.T) {
	pq := NewPriorityQueue[int](true)
	pqTestPriorityQueueDeleteImpl(t, pq, []int{15, 19, 12, 8, 13}, []int{19, 13, 12, 8}, 15)
}

func Test_PriorityQueueDeleteMaxHeap(t *testing.T) {
	pq := NewPriorityQueue[int](true)
	pqTestPriorityQueueDeleteImpl(t, pq, []int{15, 19, 12, 8, 13}, []int{19, 15, 13, 8}, 12)
}

func pqTestPriorityQueueDeleteImpl[T comparable](t *testing.T, q *PriorityQueue[T], input, expected []T, val T) {
	for i := 0; i < len(input); i++ {
		q.Add(input[i])
	}

	q.Remove(val)
	require.Equal(t, len(input)-1, q.Len())
	assert.False(t, q.Contains(val))
	for i := 0; i < len(expected); i++ {
		val, ok := q.Poll()
		assert.True(t, ok)
		assert.Equal(t, expected[i], val)
	}
	require.Zero(t, q.Len())
}
