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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_QueueLen(t *testing.T) {
	q := New[int]()
	q.Add(5)
	q.Add(6)
	assert.Equal(t, 2, q.Len())
}

func Test_QueuePeek(t *testing.T) {
	q := New[string]()
	q.Add("hello")
	q.Add("world")

	val, ok := q.Peek()
	assert.True(t, ok)
	assert.Equal(t, "hello", val)

	val, ok = q.Peek()
	assert.True(t, ok)
	assert.Equal(t, "hello", val)

	q.Poll()
	q.Poll()

	val3, ok := q.Peek()
	assert.False(t, ok)
	assert.Empty(t, val3)
}

func Test_QueueValue(t *testing.T) {
	q := New[int]()
	q.Add(15)
	q.Add(19)
	q.Add(12)
	q.Add(8)
	q.Add(1)
	q.Add(13)

	// len
	require.Equal(t, 6, q.Len())

	// Peek
	val, ok := q.Peek()
	assert.True(t, ok)
	require.Equal(t, 15, val)
	require.Equal(t, 6, q.Len())

	// Contains
	require.True(t, q.Contains(12))
	require.False(t, q.Contains(10000))

	// Poll
	val, ok = q.Poll()
	assert.True(t, ok)
	require.Equal(t, 15, val)
	require.Equal(t, 5, q.Len())

	val, ok = q.Poll()
	assert.True(t, ok)
	require.Equal(t, 19, val)
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

func Test_QueuePoll(t *testing.T) {
	q := New[string]()
	q.Add("hello")
	q.Add("world")

	val, ok := q.Poll()
	assert.True(t, ok)
	assert.Equal(t, "hello", val)

	val, ok = q.Poll()
	assert.True(t, ok)
	assert.Equal(t, "world", val)

	val, ok = q.Poll()
	assert.False(t, ok)
	assert.Empty(t, val)
}

func Test_QueueIsEmpty(t *testing.T) {
	q := New[int]()
	q.Add(5)
	q.Add(6)
	assert.False(t, q.IsEmpty())

	q.Clear()
	assert.Equal(t, 0, q.Len())
	assert.True(t, q.IsEmpty())
}

func Test_QueueInit(t *testing.T) {
	q := New[int]()
	q.Add(5)
	q.Add(6)
	q.Clear()

	assert.Equal(t, 0, q.Len())
}
