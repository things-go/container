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

func Test_QuickQueueLen(t *testing.T) {
	q := NewQuickQueue[int]()
	q.Add(5)
	q.Add(6)
	assert.Equal(t, 2, q.Len())
}

func Test_QuickQueuePeek(t *testing.T) {
	q := NewQuickQueue[string]()
	q.Add("hello")
	q.Add("world")

	val1, ok := q.Peek()
	assert.True(t, ok)
	assert.Equal(t, "hello", val1)

	val2, ok := q.Peek()
	assert.True(t, ok)
	assert.Equal(t, "hello", val2)

	q.Poll()
	q.Poll()

	val3, ok := q.Peek()
	assert.False(t, ok)
	assert.Empty(t, val3)
}

func Test_QuickQueueValue(t *testing.T) {
	q := NewQuickQueue[int]()
	q.Add(15)
	q.Add(11)
	q.Add(19)
	q.Add(12)
	q.Add(8)
	q.Add(1)

	// len
	require.Equal(t, 6, q.Len())

	// Peek
	val, ok := q.Peek()
	assert.True(t, ok)
	require.Equal(t, 15, val)
	require.Equal(t, 6, q.Len())

	// Contains
	require.True(t, q.Contains(19))
	require.False(t, q.Contains(10000))

	// Poll
	val, ok = q.Poll()
	assert.True(t, ok)
	require.Equal(t, 15, val)
	require.Empty(t, q.head[0])
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
	require.Empty(t, q.head[1])
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

func Test_QuickQueuePoll(t *testing.T) {
	q := NewQuickQueue[string]()

	q.Add("hello")
	q.Add("world")
	val1, ok := q.Poll()
	assert.True(t, ok)
	assert.Equal(t, "hello", val1)

	val2, ok := q.Poll()
	assert.True(t, ok)
	assert.Equal(t, "world", val2)

	val3, ok := q.Poll()
	assert.False(t, ok)
	assert.Empty(t, val3)
}

func Test_QuickQueueIsEmpty(t *testing.T) {
	q := NewQuickQueue[int]()
	q.Add(5)
	q.Add(6)
	assert.False(t, q.IsEmpty())

	val, ok := q.Poll()
	assert.True(t, ok)
	assert.Equal(t, 5, val)

	val, ok = q.Peek()
	assert.True(t, ok)
	assert.Equal(t, 6, val)
	q.Clear()
	assert.Equal(t, 0, q.Len())
	assert.True(t, q.IsEmpty())
}

func Test_QuickQueueInit(t *testing.T) {
	q := NewQuickQueue[int]()
	q.Add(5)
	q.Add(6)
	q.Clear()

	assert.Equal(t, 0, q.Len())
}
