package comparator

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/things-go/container"
)

func TestSort(t *testing.T) {
	input1 := []interface{}{6, 4, 9, 19, 15}
	expected1 := []interface{}{4, 6, 9, 15, 19}
	assertSort(t, input1, expected1, false, nil)

	input2 := []interface{}{"benjamin", "alice", "john", "tom", "roy"}
	expected2 := []interface{}{"alice", "benjamin", "john", "roy", "tom"}
	assertSort(t, input2, expected2, false, nil)
}

func TestSortWithComparator(t *testing.T) {
	input1 := []interface{}{6, 4, 9, 19, 15}
	expected1 := []interface{}{19, 15, 9, 6, 4}
	assertSort(t, input1, expected1, false, CompareReverseInt)

	input2 := []interface{}{"benjamin", "alice", "john", "tom", "roy"}
	expected2 := []interface{}{"tom", "roy", "john", "benjamin", "alice"}
	assertSort(t, input2, expected2, false, CompareReverseString)
}

func TestReverseSort(t *testing.T) {
	input1 := []interface{}{6, 4, 9, 19, 15}
	expected1 := []interface{}{19, 15, 9, 6, 4}
	assertSort(t, input1, expected1, true, nil)

	input2 := []interface{}{"benjamin", "alice", "john", "tom", "roy"}
	expected2 := []interface{}{"tom", "roy", "john", "benjamin", "alice"}
	assertSort(t, input2, expected2, true, nil)
}

func TestReverseSortWithComparator(t *testing.T) {
	input1 := []interface{}{6, 4, 9, 19, 15}
	expected1 := []interface{}{4, 6, 9, 15, 19}
	assertSort(t, input1, expected1, true, CompareReverseInt)

	input2 := []interface{}{"benjamin", "alice", "john", "tom", "roy"}
	expected2 := []interface{}{"alice", "benjamin", "john", "roy", "tom"}
	assertSort(t, input2, expected2, true, CompareReverseString)
}

func assertSort(t *testing.T, input, expected []interface{}, reverse bool, c container.Comparator) {
	// sort
	contain := NewContainer(input, c)
	if reverse {
		contain.Reverse()
	}
	contain.Sort()

	for i := 0; i < len(input); i++ {
		assert.Equal(t, expected[i], input[i])
	}
}

// Compare returns Reverse order for string.
func CompareReverseString(v1, v2 interface{}) int {
	i1, i2 := v1.(string), v2.(string)

	if i1 < i2 {
		return 1
	}
	if i1 > i2 {
		return -1
	}
	return 0
}

// Compare returns Reverse order for int.
func CompareReverseInt(v1, v2 interface{}) int {
	i1, i2 := v1.(int), v2.(int)

	if i1 < i2 {
		return 1
	}
	if i1 > i2 {
		return -1
	}
	return 0
}

/***************************************heap*************************************/

func TestHeap(t *testing.T) {
	input1 := []interface{}{6, 4, 9, 19, 15}
	expected1 := []interface{}{4, 6, 9, 15, 19}
	heapTestImpl(t, input1, expected1, true, nil)

	input2 := []interface{}{"benjamin", "alice", "john", "tom", "roy"}
	expected2 := []interface{}{"alice", "benjamin", "john", "roy", "tom"}
	heapTestImpl(t, input2, expected2, true, nil)
}

func TestHeapWithComparator(t *testing.T) {
	input1 := []interface{}{6, 4, 9, 19, 15}
	expected1 := []interface{}{19, 15, 9, 6, 4}
	heapTestImpl(t, input1, expected1, true, CompareReverseInt)

	input2 := []interface{}{"benjamin", "alice", "john", "tom", "roy"}
	expected2 := []interface{}{"tom", "roy", "john", "benjamin", "alice"}
	heapTestImpl(t, input2, expected2, true, CompareReverseString)
}

func TestMaxHeap(t *testing.T) {
	input1 := []interface{}{6, 4, 9, 19, 15}
	expected1 := []interface{}{19, 15, 9, 6, 4}
	heapTestImpl(t, input1, expected1, false, nil)

	input2 := []interface{}{"benjamin", "alice", "john", "tom", "roy"}
	expected2 := []interface{}{"tom", "roy", "john", "benjamin", "alice"}
	heapTestImpl(t, input2, expected2, false, nil)
}

func TestMaxHeapWithComparator(t *testing.T) {
	input1 := []interface{}{6, 4, 9, 19, 15}
	expected1 := []interface{}{4, 6, 9, 15, 19}
	heapTestImpl(t, input1, expected1, false, CompareReverseInt)

	input2 := []interface{}{"benjamin", "alice", "john", "tom", "roy"}
	expected2 := []interface{}{"alice", "benjamin", "john", "roy", "tom"}
	heapTestImpl(t, input2, expected2, false, CompareReverseString)
}

func heapTestImpl(t *testing.T, input, expected []interface{}, isMinHeap bool, c container.Comparator) {
	contain := &Container{
		Items:   input,
		compare: c,
		reverse: !isMinHeap,
	}
	heap.Init(contain)

	// Pop all elements from heap
	for i := 0; i < len(expected); i++ {
		require.Equal(t, expected[i], heap.Pop(contain))
	}
	assert.Zero(t, contain.Len())
}

func TestHeapRemove(t *testing.T) {
	input1 := []interface{}{6, 4, 9, 19, 15}
	expected1 := []interface{}{4, 6, 9, 19}
	heapRemoveTestImpl(t, input1, expected1, 15, true, nil)

	input2 := []interface{}{"benjamin", "alice", "john", "tom", "roy"}
	expected2 := []interface{}{"alice", "benjamin", "roy", "tom"}
	heapRemoveTestImpl(t, input2, expected2, "john", true, nil)
}

func TestHeapRemoveWithComparator(t *testing.T) {
	input1 := []interface{}{6, 4, 9, 19, 15}
	expected1 := []interface{}{19, 15, 6, 4}
	heapRemoveTestImpl(t, input1, expected1, 9, true, CompareReverseInt)

	input2 := []interface{}{"benjamin", "alice", "john", "tom", "roy"}
	expected2 := []interface{}{"tom", "john", "benjamin", "alice"}
	heapRemoveTestImpl(t, input2, expected2, "roy", true, CompareReverseString)
}

func TestMaxHeapRemove(t *testing.T) {
	input1 := []interface{}{6, 4, 9, 19, 15}
	expected1 := []interface{}{15, 9, 6, 4}
	heapRemoveTestImpl(t, input1, expected1, 19, false, nil)

	input2 := []interface{}{"benjamin", "alice", "john", "tom", "roy"}
	expected2 := []interface{}{"tom", "roy", "john", "benjamin"}
	heapRemoveTestImpl(t, input2, expected2, "alice", false, nil)
}

func TestMaxHeapRemoveWithComparator(t *testing.T) {
	input1 := []interface{}{6, 4, 9, 19, 15}
	expected1 := []interface{}{4, 6, 9, 15}
	heapRemoveTestImpl(t, input1, expected1, 19, false, CompareReverseInt)

	input2 := []interface{}{"benjamin", "alice", "john", "tom", "roy"}
	expected2 := []interface{}{"alice", "benjamin", "john", "roy"}
	heapRemoveTestImpl(t, input2, expected2, "tom", false, CompareReverseString)
}

func heapRemoveTestImpl(t *testing.T, input, expected []interface{},
	val interface{}, isMinHeap bool, c container.Comparator) {
	contain := &Container{
		Items:   input,
		compare: c,
		reverse: !isMinHeap,
	}
	heap.Init(contain)

	// find the index of the value to be removed
	index := 0
	for i := 0; i < len(input); i++ {
		if input[i] == val {
			index = i
			break
		}
	}

	// call HeapPreRemove
	v := heap.Remove(contain, index)
	require.Equal(t, v, val)
	require.Equal(t, nil, input[len(input)-1])

	// Pop all elements from heap one by one
	for i := 0; i < len(expected); i++ {
		require.Equal(t, expected[i], heap.Pop(contain))
	}
	assert.Zero(t, contain.Len())
}

// Test: HeapInit and HeapPostUpdate.
func TestHeapFix(t *testing.T) {
	input1 := []interface{}{6, 4, 9, 19, 15}
	expected1 := []interface{}{4, 6, 9, 19, 25}
	heapFixTestImpl(t, input1, expected1, 15, 25, true, nil)

	input2 := []interface{}{"benjamin", "alice", "john", "tom", "roy"}
	expected2 := []interface{}{"alice", "benjamin", "ken", "roy", "tom"}
	heapFixTestImpl(t, input2, expected2, "john", "ken", true, nil)
}

func TestHeapFixWithComparator(t *testing.T) {
	input1 := []interface{}{6, 4, 9, 19, 15}
	expected1 := []interface{}{19, 15, 6, 4, 3}
	heapFixTestImpl(t, input1, expected1, 9, 3, true, CompareReverseInt)

	input2 := []interface{}{"benjamin", "alice", "john", "tom", "roy"}
	expected2 := []interface{}{"tom", "john", "benjamin", "alice", "ali"}
	heapFixTestImpl(t, input2, expected2, "roy", "ali", true, CompareReverseString)
}

func TestMaxHeapFix(t *testing.T) {
	input1 := []interface{}{6, 4, 9, 19, 15}
	expected1 := []interface{}{15, 13, 9, 6, 4}
	heapFixTestImpl(t, input1, expected1, 19, 13, false, nil)

	input2 := []interface{}{"benjamin", "alice", "john", "tom", "roy"}
	expected2 := []interface{}{"trevor", "tom", "roy", "john", "benjamin"}
	heapFixTestImpl(t, input2, expected2, "alice", "trevor", false, nil)
}

func TestMaxHeapFixWithComparator(t *testing.T) {
	input1 := []interface{}{6, 4, 9, 19, 15}
	expected1 := []interface{}{4, 6, 7, 9, 15}
	heapFixTestImpl(t, input1, expected1, 19, 7, false, CompareReverseInt)

	input2 := []interface{}{"benjamin", "alice", "john", "tom", "roy"}
	expected2 := []interface{}{"alice", "benjamin", "john", "roy", "zoo"}
	heapFixTestImpl(t, input2, expected2, "tom", "zoo", false, CompareReverseString)
}

func heapFixTestImpl(t *testing.T, input, expected []interface{},
	oldVal, newVal interface{}, isMinHeap bool, c container.Comparator) {
	in := &Container{
		Items:   input,
		compare: c,
		reverse: !isMinHeap,
	}
	heap.Init(in)
	// find the index of the value to be updated
	index := 0
	for i := 0; i < len(input); i++ {
		if input[i] == oldVal {
			index = i
			// update the value
			input[index] = newVal
			break
		}
	}
	heap.Fix(in, index)

	// Pop all elements from heap one by one
	for i := 0; i < len(expected); i++ {
		require.Equal(t, expected[i], heap.Pop(in))
	}

	assert.Zero(t, in.Len())

	heap.Push(in, 1)
	assert.Equal(t, 1, heap.Pop(in))

	// improve couver
	in.Sort()
}
