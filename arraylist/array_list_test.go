package arraylist

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ArrayListLen(t *testing.T) {
	l := New[int]()

	l.PushBack(5)
	l.PushBack(6)
	l.PushBack(7)
	assert.Equal(t, 3, l.Len())

	// remove the element at the position 1
	v, err := l.Remove(1)
	assert.Nil(t, err)
	assert.Equal(t, 6, v)
	assert.Equal(t, 2, l.Len())
	assert.False(t, l.IsEmpty())
	assert.False(t, l.Contains(6))

	v, err = l.Remove(100)
	assert.NotNil(t, err)
	assert.Empty(t, v)

	// clear l the elements
	l.Clear()
	assert.True(t, l.IsEmpty())
}

func Test_ArrayListValue(t *testing.T) {
	l := New[int]()
	l.Push(5)
	l.PushBack(7)
	l.PushFront(6)

	require.True(t, slices.Equal(l.Values(), []int{6, 5, 7}))
	// peek
	val, ok := l.Peek()
	assert.True(t, ok)
	assert.Equal(t, 6, val)

	val, ok = l.PeekFront()
	assert.True(t, ok)
	assert.Equal(t, 6, val)

	val, ok = l.PeekBack()
	assert.True(t, ok)
	assert.Equal(t, 7, val)

	err := l.Add(2, 8)
	assert.Nil(t, err)

	require.True(t, slices.Equal(l.Values(), []int{6, 5, 8, 7}))

	v, err := l.Get(2)
	assert.Nil(t, err)
	assert.Equal(t, 8, v)

	v, err = l.Get(3)
	assert.Nil(t, err)
	assert.Equal(t, 7, v)

	// check an element which doesn't exist
	assert.False(t, l.Contains(9))
	assert.False(t, l.RemoveValue(9))

	// check element 8
	assert.True(t, l.Contains(8))
	assert.True(t, l.RemoveValue(8))
	assert.False(t, l.Contains(8))

	require.True(t, slices.Equal(l.Values(), []int{6, 5, 7}))

	// get out of range
	v, err = l.Get(l.Len())
	assert.NotNil(t, err)
	assert.Empty(t, v)
	v, err = l.Get(-1)
	assert.NotNil(t, err)
	assert.Empty(t, v)

	// check length at last
	assert.Equal(t, 3, l.Len())

	val, ok = l.Poll()
	assert.True(t, ok)
	assert.Equal(t, 6, val)

	val, ok = l.PollBack()
	assert.True(t, ok)
	assert.Equal(t, 7, val)

	val, ok = l.PollBack()
	assert.True(t, ok)
	assert.Equal(t, 5, val)

	require.True(t, l.IsEmpty())

	val, ok = l.PollFront()
	assert.False(t, ok)
	assert.Empty(t, val)
	val, ok = l.PollBack()
	assert.False(t, ok)
	assert.Empty(t, val)

	l.Clear()
	val, ok = l.Peek()
	assert.False(t, ok)
	assert.Empty(t, val)
	val, ok = l.PeekFront()
	assert.False(t, ok)
	assert.Empty(t, val)
	val, ok = l.PeekBack()
	assert.False(t, ok)
	assert.Empty(t, val)

	// nothing remove
	assert.False(t, l.RemoveValue(8))
	err = l.Add(0, 1)
	assert.Nil(t, err)

	require.True(t, slices.Equal(l.Values(), []int{1}))

	// invalid index
	err = l.Add(-1, 1)
	assert.NotNil(t, err)
	err = l.Add(l.Len()+1, 1)
	assert.Error(t, err)
}

func Test_ArrayListIterator(t *testing.T) {
	l := New[int]()
	items := []int{5, 6, 7}
	l.PushBack(5)
	l.PushBack(6)
	l.PushBack(7)
	idx := 0
	l.Iterator(func(v int) bool {
		assert.Equal(t, items[idx], v)
		idx++
		return true
	})
	l.Iterator(nil)
}

func Test_ArrayListReverseIterator(t *testing.T) {
	items := []int{5, 6, 7}
	l := New[int]()
	l.PushBack(5)
	l.PushBack(6)
	l.PushBack(7)
	idx := len(items) - 1
	l.ReverseIterator(func(v int) bool {
		assert.Equal(t, items[idx], v)
		idx--
		return true
	})
	l.ReverseIterator(nil)
}

func Test_ArrayListSort(t *testing.T) {
	ll := New[int]()

	expect := []int{4, 6, 7, 15}

	ll.PushBack(15)
	ll.PushBack(6)
	ll.PushBack(7)
	ll.PushBack(4)

	// sort
	ll.Sort(func(i, j int) int { return i - j })
	assert.Equal(t, 4, ll.Len())
	for i := 0; i < ll.Len(); i++ {
		v, err := ll.Get(i)
		assert.Nil(t, err)
		assert.Equal(t, expect[i], v)
	}

	// reverse sorting
	ll.Sort(func(i, j int) int { return j - i })
	assert.Equal(t, 4, ll.Len())
	for i := 0; i < ll.Len(); i++ {
		v, err := ll.Get(i)
		assert.Nil(t, err)
		assert.Equal(t, expect[ll.Len()-1-i], v)
	}
}

func Test_Extending(t *testing.T) {
	l1 := New[int]()
	l2 := New[int]()

	l1.PushBack(1)
	l1.PushBack(2)
	l1.PushBack(3)

	l2.PushBack(4)
	l2.PushBack(5)

	l3 := New[int]()
	l3.PushBackList(l1)
	require.True(t, slices.Equal(l3.Values(), []int{1, 2, 3}))
	l3.PushBackList(l2)
	require.True(t, slices.Equal(l3.Values(), []int{1, 2, 3, 4, 5}))

	l3 = New[int]()
	l3.PushFrontList(l2)
	require.True(t, slices.Equal(l3.Values(), []int{4, 5}))
	l3.PushFrontList(l1)
	require.True(t, slices.Equal(l3.Values(), []int{1, 2, 3, 4, 5}))

	require.True(t, slices.Equal(l1.Values(), []int{1, 2, 3}))
	require.True(t, slices.Equal(l2.Values(), []int{4, 5}))

	l3 = New[int]()
	l3.PushBackList(l1)
	require.True(t, slices.Equal(l3.Values(), []int{1, 2, 3}))
	l3.PushBackList(l3)
	require.True(t, slices.Equal(l3.Values(), []int{1, 2, 3, 1, 2, 3}))

	l3 = New[int]()
	l3.PushFrontList(l1)
	require.True(t, slices.Equal(l3.Values(), []int{1, 2, 3}))
	l3.PushFrontList(l3)
	require.True(t, slices.Equal(l3.Values(), []int{1, 2, 3, 1, 2, 3}))

	l3 = New[int]()
	l1.PushBackList(l3)
	require.True(t, slices.Equal(l1.Values(), []int{1, 2, 3}))
	l1.PushFrontList(l3)
	require.True(t, slices.Equal(l1.Values(), []int{1, 2, 3}))
}
