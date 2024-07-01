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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LinkedMapLen(t *testing.T) {
	lm := New[int, string]()
	lm.Push(24, "benjamin")
	lm.Push(43, "alice")
	lm.Push(18, "john")
	assert.Equal(t, 0, lm.capacity)
	assert.Equal(t, 3, lm.Len())

	// exist
	v, ok := lm.Remove(43)
	require.True(t, ok)
	require.Equal(t, "alice", v)
	require.Equal(t, 2, lm.Len())
	// not exist
	v, ok = lm.Remove(1000)
	require.False(t, ok)
	require.Empty(t, v)
	require.Equal(t, 2, lm.Len())

	lm.Clear()
	require.Equal(t, 0, lm.Len())
	require.True(t, lm.IsEmpty())

	// over capacity
	lm = New[int, string](WithCap[int, string](3))

	// not exist
	k, v, exist := lm.Peek()
	require.False(t, exist)
	assert.Empty(t, k)
	assert.Empty(t, v)

	k, v, exist = lm.PeekBack()
	require.False(t, exist)
	assert.Empty(t, k)
	assert.Empty(t, v)

	lm.Push(24, "benjamin")
	lm.Push(43, "alice")
	lm.Push(18, "john")
	lm.Push(25, "haha")

	assert.Equal(t, 3, lm.Cap())
	assert.Equal(t, 3, lm.Len())

	k, v, exist = lm.Peek()
	require.True(t, exist)
	assert.Equal(t, 43, k)
	assert.Equal(t, "alice", v)

	// exist
	lm.Push(18, "john")
	k, v, exist = lm.PeekBack()
	require.True(t, exist)
	assert.Equal(t, 18, k)
	assert.Equal(t, "john", v)

	// exist
	lm.PushFront(25, "haha")
	k, v, exist = lm.Peek()
	require.True(t, exist)
	assert.Equal(t, 25, k)
	assert.Equal(t, "haha", v)

	// not exist
	lm.PushFront(99, "noexist")
	k, v, exist = lm.Peek()
	require.True(t, exist)
	assert.Equal(t, 99, k)
	assert.Equal(t, "noexist", v)

	// last check len and capacity
	assert.Equal(t, 3, lm.Cap())
	assert.Equal(t, 3, lm.Len())
}

func Test_LinkedMapValue(t *testing.T) {
	equal := func(a, b string) bool { return a == b }
	lm := New[int, string]()
	keys := []int{24, 43, 18, 23, 35}
	values := []string{"benjamin", "alice", "john", "tom", "bill"}
	for i := 0; i < len(keys); i++ {
		lm.Push(keys[i], values[i])
	}
	// not found
	assert.False(t, lm.ContainsValue("haha", equal))

	// test Contains & ContainsValue
	for _, k := range keys {
		assert.True(t, lm.Contains(k))
	}

	for _, v := range values {
		assert.True(t, lm.ContainsValue(v, equal))
	}

	// not found
	assert.Empty(t, lm.Get(1000))

	// test Get & GetOrDefault
	for i, k := range keys {
		v := lm.Get(k)
		assert.Equal(t, values[i], v)
	}

	v := lm.Get(50, "defaultName")
	assert.Equal(t, "defaultName", v)

	// test Remove, Poll and PollBack
	v, ok := lm.Remove(43)
	assert.False(t, !ok || v != "alice")

	k, v, ok := lm.Poll()
	assert.False(t, k != 24 || v != "benjamin" || !ok)

	k, v, ok = lm.PollBack()
	assert.False(t, k != 35 || v != "bill" || !ok)
	assert.Equal(t, 2, lm.Len())

	// not exist
	lm.Clear()
	k, v, ok = lm.Poll()
	assert.False(t, ok)
	assert.Empty(t, k)
	assert.Empty(t, v)

	k, v, ok = lm.PollBack()
	assert.False(t, ok)
	assert.Empty(t, k)
	assert.Empty(t, v)
}

func Test_LinkedMapIterate(t *testing.T) {
	lm := New[int, string]()
	keys := []int{24, 43, 18, 23, 35}
	values := []string{"benjamin", "alice", "john", "tom", "bill"}
	for i := 0; i < len(keys); i++ {
		lm.Push(keys[i], values[i])
	}

	idx := 0
	lm.Iterator(func(k int, v string) bool {
		assert.Equal(t, keys[idx], k)
		assert.Equal(t, values[idx], v)
		idx++
		return true
	})
	idx = lm.Len() - 1
	lm.ReverseIterator(func(k int, v string) bool {
		assert.Equal(t, keys[idx], k)
		assert.Equal(t, values[idx], v)
		idx--
		return true
	})

	// improve cover
	lm.Iterator(nil)
	lm.ReverseIterator(nil)
}
