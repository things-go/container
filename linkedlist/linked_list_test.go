package linkedlist

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkList(t *testing.T, l *LinkedList, es []interface{}) {
	require.Equal(t, len(es), l.Len())
	for i, e := 0, l.l.Front(); e != nil; e = e.Next() {
		assert.Equal(t, es[i], e.Value.(int))
		i++
	}
}

func TestLinkedListLen(t *testing.T) {
	l := New()

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

	v, err = l.Remove(100)
	assert.NotNil(t, err)
	assert.Nil(t, v)

	// clear l the elements
	l.Clear()
	assert.Equal(t, 0, l.Len())
	assert.True(t, l.IsEmpty())
}

func TestLinkedListValue(t *testing.T) {
	l := New()
	l.Push(5)
	l.PushBack(7)
	l.PushFront(6)

	assert.Equal(t, 6, l.Peek())
	assert.Equal(t, 6, l.PeekFront())
	assert.Equal(t, 7, l.PeekBack())

	err := l.Add(2, 8)
	assert.Nil(t, err)

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
	assert.False(t, l.Contains(9))
	assert.True(t, l.RemoveValue(8))
	assert.False(t, l.Contains(9))

	// get out of range
	v, err = l.Get(l.Len())
	assert.NotNil(t, err)
	assert.Nil(t, v)
	v, err = l.Get(-1)
	assert.NotNil(t, err)
	assert.Nil(t, v)

	// check length at last
	assert.Equal(t, 3, l.Len())

	assert.Equal(t, 6, l.Poll())
	assert.Equal(t, 7, l.PollBack())
	assert.Equal(t, 5, l.PollBack())

	assert.Nil(t, l.PollFront())
	assert.Nil(t, l.PollBack())

	l.Clear()
	assert.Nil(t, l.Peek())
	assert.Nil(t, l.PeekFront())
	assert.Nil(t, l.PeekBack())

	// nothing remove
	assert.False(t, l.RemoveValue(8))
	err = l.Add(0, 1)
	assert.Nil(t, err)

	// invalid index
	err = l.Add(-1, 1)
	assert.NotNil(t, err)
	err = l.Add(l.Len()+1, 1)
	assert.NotNil(t, err)
}

func TestUserCompare(t *testing.T) {
	ll := New(WithComparator(&linkedListNode{}))
	ll.PushBack(&linkedListNode{age: 32})
	ll.PushBack(&linkedListNode{age: 20})
	ll.PushBack(&linkedListNode{age: 27})
	ll.PushBack(&linkedListNode{age: 25})

	idx := ll.indexOf(&linkedListNode{age: 20})
	assert.Equal(t, 1, idx)

	ok := ll.RemoveValue(&linkedListNode{age: 20})
	assert.True(t, ok)
	assert.Equal(t, 3, ll.Len())
}

func TestLinkedListIterator(t *testing.T) {
	l := New()
	items := []int{5, 6, 7}
	l.PushBack(5)
	l.PushBack(6)
	l.PushBack(7)
	idx := 0
	l.Iterator(func(v interface{}) bool {
		assert.Equal(t, items[idx], v)
		idx++
		return true
	})
	l.Iterator(nil)
}

func TestLinkedListReverseIterator(t *testing.T) {
	items := []int{5, 6, 7}
	l := New()
	l.PushBack(5)
	l.PushBack(6)
	l.PushBack(7)
	idx := len(items) - 1
	l.ReverseIterator(func(v interface{}) bool {
		assert.Equal(t, items[idx], v)
		idx--
		return true
	})
	l.ReverseIterator(nil)
}

func TestLinkedListSort(t *testing.T) {
	ll := New()

	expect := []int{4, 6, 7, 15}

	ll.PushBack(15)
	ll.PushBack(6)
	ll.PushBack(7)
	ll.PushBack(4)

	// sort
	ll.Sort()
	assert.Equal(t, 4, ll.Len())
	for i := 0; i < ll.Len(); i++ {
		v, err := ll.Get(i)
		assert.Nil(t, err)
		assert.Equal(t, expect[i], v)
	}

	// reverse sorting
	ll.Sort(true)
	assert.Equal(t, 4, ll.Len())
	for i := 0; i < ll.Len(); i++ {
		v, err := ll.Get(i)
		assert.Nil(t, err)
		assert.Equal(t, expect[ll.Len()-1-i], v)
	}
}

func TestExtending(t *testing.T) {
	l1 := New()
	l2 := New()

	l1.PushBack(1)
	l1.PushBack(2)
	l1.PushBack(3)

	l2.PushBack(4)
	l2.PushBack(5)

	l3 := New()
	l3.PushBackList(l1)
	checkList(t, l3, []interface{}{1, 2, 3})
	l3.PushBackList(l2)
	checkList(t, l3, []interface{}{1, 2, 3, 4, 5})

	l3 = New()
	l3.PushFrontList(l2)
	checkList(t, l3, []interface{}{4, 5})
	l3.PushFrontList(l1)
	checkList(t, l3, []interface{}{1, 2, 3, 4, 5})

	checkList(t, l1, []interface{}{1, 2, 3})
	checkList(t, l2, []interface{}{4, 5})

	l3 = New()
	l3.PushBackList(l1)
	checkList(t, l3, []interface{}{1, 2, 3})
	l3.PushBackList(l3)
	checkList(t, l3, []interface{}{1, 2, 3, 1, 2, 3})

	l3 = New()
	l3.PushFrontList(l1)
	checkList(t, l3, []interface{}{1, 2, 3})
	l3.PushFrontList(l3)
	checkList(t, l3, []interface{}{1, 2, 3, 1, 2, 3})

	l3 = New()
	l1.PushBackList(l3)
	checkList(t, l1, []interface{}{1, 2, 3})
	l1.PushFrontList(l3)
	checkList(t, l1, []interface{}{1, 2, 3})
}

func TestLinkdedListComparatorSort(t *testing.T) {
	expect := []*linkedListNode{{age: 20}, {age: 25}, {age: 27}, {age: 32}}
	ll := New(WithComparator(&linkedListNode{}))
	ll.PushBack(&linkedListNode{age: 32})
	ll.PushBack(&linkedListNode{age: 20})
	ll.PushBack(&linkedListNode{age: 27})
	ll.PushBack(&linkedListNode{age: 25})

	// sort
	ll.Sort()
	assert.Equal(t, 4, ll.Len())
	for i := 0; i < ll.Len(); i++ {
		v, err := ll.Get(i)
		assert.Nil(t, err)
		assert.Equal(t, expect[i], v)
	}

	// reverse sorting
	ll.Sort(true)
	assert.Equal(t, 4, ll.Len())
	for i := 0; i < ll.Len(); i++ {
		v, err := ll.Get(i)
		assert.Nil(t, err)
		assert.Equal(t, expect[ll.Len()-1-i], v)
	}

	ll.Clear()
	ll.Sort()

	value := ll.Values()
	assert.Equal(t, []interface{}{}, value)
}

type linkedListNode struct {
	age int
}

func (aln *linkedListNode) Compare(v1, v2 interface{}) int {
	n1, n2 := v1.(*linkedListNode), v2.(*linkedListNode)

	if n1.age < n2.age {
		return -1
	}

	if n1.age == n2.age {
		return 0
	}

	return 1
}
