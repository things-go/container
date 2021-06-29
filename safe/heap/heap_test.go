package heap

import (
	"sync"
	"testing"
	"time"
)

func testHeapObjectKeyFunc(obj interface{}) (string, error) {
	return obj.(testHeapObject).name, nil
}

type testHeapObject struct {
	name string
	val  interface{}
}

func mkHeapObj(name string, val interface{}) testHeapObject {
	return testHeapObject{name: name, val: val}
}

func compareInts(val1, val2 interface{}) int {
	first := val1.(testHeapObject).val.(int)
	second := val2.(testHeapObject).val.(int)
	if first < second {
		return -1
	}
	if first > second {
		return 1
	}
	return 0
}

// TestHeapBasic tests Heap invariant and synchronization.
func TestHeapBasic(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	var wg sync.WaitGroup
	wg.Add(2)
	const amount = 500
	var i, u int
	// Insert items in the heap in opposite orders in two go routines.
	go func() {
		for i = amount; i > 0; i-- {
			h.Add(mkHeapObj(string([]rune{'a', rune(i)}), i)) // nolint: errcheck
		}
		wg.Done()
	}()
	go func() {
		for u = 0; u < amount; u++ {
			h.Add(mkHeapObj(string([]rune{'b', rune(u)}), u+1)) // nolint: errcheck
		}
		wg.Done()
	}()
	// Wait for the two go routines to finish.
	wg.Wait()
	// Make sure that the numbers are popped in ascending order.
	prevNum := 0
	for i := 0; i < amount*2; i++ {
		obj, err := h.Pop()
		num := obj.(testHeapObject).val.(int)
		// All the items must be sorted.
		if err != nil || prevNum > num {
			t.Errorf("got %v out of order, last was %v", obj, prevNum)
		}
		prevNum = num
	}
}

// Tests Heap.Push and ensures that heap invariant is preserved after adding items.
func TestHeap_Add(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	h.Add(mkHeapObj("foo", 10)) // nolint: errcheck
	h.Add(mkHeapObj("bar", 1))  // nolint: errcheck
	h.Add(mkHeapObj("baz", 11)) // nolint: errcheck
	h.Add(mkHeapObj("zab", 30)) // nolint: errcheck
	h.Add(mkHeapObj("foo", 13)) // nolint: errcheck

	item, err := h.Pop()
	if e, a := 1, item.(testHeapObject).val; err != nil || a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	item, err = h.Pop()
	if e, a := 11, item.(testHeapObject).val; err != nil || a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	// Nothing is deleted.
	h.Delete(mkHeapObj("baz", 11)) // nolint: errcheck
	// foo is updated.
	h.Add(mkHeapObj("foo", 14)) // nolint: errcheck
	item, err = h.Pop()
	if e, a := 14, item.(testHeapObject).val; err != nil || a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	item, err = h.Pop()
	if e, a := 30, item.(testHeapObject).val; err != nil || a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
}

// TestHeap_BulkAdd tests Heap.BulkAdd functionality and ensures that all the
// items given to BulkAdd are added to the queue before Pop reads them.
func TestHeap_BulkAdd(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	const amount = 500
	// Insert items in the heap in opposite orders in a go routine.
	go func() {
		l := []interface{}{}
		for i := amount; i > 0; i-- {
			l = append(l, mkHeapObj(string([]rune{'a', rune(i)}), i))
		}
		h.BulkAdd(l) // nolint: errcheck
	}()
	prevNum := -1
	for i := 0; i < amount; i++ {
		obj, err := h.Pop()
		num := obj.(testHeapObject).val.(int)
		// All the items must be sorted.
		if err != nil || prevNum >= num {
			t.Errorf("got %v out of order, last was %v", obj, prevNum)
		}
		prevNum = num
	}
}

// TestHeapEmptyPop tests that pop returns properly after heap is closed.
func TestHeapEmptyPop(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	go func() {
		time.Sleep(1 * time.Second)
		h.Close()
	}()
	_, err := h.Pop()
	if err == nil || err != ErrHeapClosed {
		t.Errorf("pop should have returned heap closed error: %v", err)
	}
}

// TestHeap_AddIfNotPresent tests Heap.AddIfNotPresent and ensures that heap
// invariant is preserved after adding items.
func TestHeap_AddIfNotPresent(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	h.AddIfNotPresent(mkHeapObj("foo", 10)) // nolint: errcheck
	h.AddIfNotPresent(mkHeapObj("bar", 1))  // nolint: errcheck
	h.AddIfNotPresent(mkHeapObj("baz", 11)) // nolint: errcheck
	h.AddIfNotPresent(mkHeapObj("zab", 30)) // nolint: errcheck
	// This is not added.
	h.AddIfNotPresent(mkHeapObj("foo", 13)) // nolint: errcheck

	if length := len(h.data.items); length != 4 {
		t.Errorf("unexpected number of items: %d", length)
	}
	if val := h.data.items["foo"].obj.(testHeapObject).val; val != 10 {
		t.Errorf("unexpected value: %d", val)
	}
	item, err := h.Pop()
	if e, a := 1, item.(testHeapObject).val; err != nil || a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	item, err = h.Pop()
	if e, a := 10, item.(testHeapObject).val; err != nil || a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	// bar is already popped. Let's add another one.
	h.AddIfNotPresent(mkHeapObj("bar", 14)) // nolint: errcheck
	item, err = h.Pop()
	if e, a := 11, item.(testHeapObject).val; err != nil || a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	item, err = h.Pop()
	if e, a := 14, item.(testHeapObject).val; err != nil || a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
}

// TestHeap_Delete tests Heap.Delete and ensures that heap invariant is
// preserved after deleting items.
func TestHeap_Delete(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	h.Add(mkHeapObj("foo", 10)) // nolint: errcheck
	h.Add(mkHeapObj("bar", 1))  // nolint: errcheck
	h.Add(mkHeapObj("bal", 31)) // nolint: errcheck
	h.Add(mkHeapObj("baz", 11)) // nolint: errcheck

	// Delete head. Delete should work with "key" and doesn't care about the value.
	if err := h.Delete(mkHeapObj("bar", 200)); err != nil {
		t.Fatalf("Failed to delete head.")
	}
	item, err := h.Pop()
	if e, a := 10, item.(testHeapObject).val; err != nil || a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	h.Add(mkHeapObj("zab", 30)) // nolint: errcheck
	h.Add(mkHeapObj("faz", 30)) // nolint: errcheck
	length := h.data.Len()
	// Delete non-existing item.
	if err = h.Delete(mkHeapObj("non-existent", 10)); err == nil || length != h.data.Len() {
		t.Fatalf("Didn't expect any item removal")
	}
	// Delete tail.
	if err = h.Delete(mkHeapObj("bal", 31)); err != nil {
		t.Fatalf("Failed to delete tail.")
	}
	// Delete one of the items with value 30.
	if err = h.Delete(mkHeapObj("zab", 30)); err != nil {
		t.Fatalf("Failed to delete item.")
	}
	item, err = h.Pop()
	if e, a := 11, item.(testHeapObject).val; err != nil || a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	item, err = h.Pop()
	if e, a := 30, item.(testHeapObject).val; err != nil || a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	if h.data.Len() != 0 {
		t.Fatalf("expected an empty heap.")
	}
}

// TestHeap_Update tests Heap.Update and ensures that heap invariant is
// preserved after adding items.
func TestHeap_Update(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	h.Add(mkHeapObj("foo", 10)) // nolint: errcheck
	h.Add(mkHeapObj("bar", 1))  // nolint: errcheck
	h.Add(mkHeapObj("bal", 31)) // nolint: errcheck
	h.Add(mkHeapObj("baz", 11)) // nolint: errcheck

	// Update an item to a value that should push it to the head.
	h.Update(mkHeapObj("baz", 0)) // nolint: errcheck
	if h.data.queue[0] != "baz" || h.data.items["baz"].index != 0 {
		t.Fatalf("expected baz to be at the head")
	}
	item, err := h.Pop()
	if e, a := 0, item.(testHeapObject).val; err != nil || a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	// Update bar to push it farther back in the queue.
	h.Update(mkHeapObj("bar", 100)) // nolint: errcheck
	if h.data.queue[0] != "foo" || h.data.items["foo"].index != 0 {
		t.Fatalf("expected foo to be at the head")
	}
}

// TestHeap_Get tests Heap.Get.
func TestHeap_Get(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	h.Add(mkHeapObj("foo", 10)) // nolint: errcheck
	h.Add(mkHeapObj("bar", 1))  // nolint: errcheck
	h.Add(mkHeapObj("bal", 31)) // nolint: errcheck
	h.Add(mkHeapObj("baz", 11)) // nolint: errcheck

	// Get works with the key.
	obj, exists, err := h.Get(mkHeapObj("baz", 0))
	if err != nil || exists == false || obj.(testHeapObject).val != 11 {
		t.Fatalf("unexpected error in getting element")
	}
	// Get non-existing object.
	_, exists, err = h.Get(mkHeapObj("non-existing", 0))
	if err != nil || exists == true {
		t.Fatalf("didn't expect to get any object")
	}
}

// TestHeap_GetByKey tests Heap.GetByKey and is very similar to TestHeap_Get.
func TestHeap_GetByKey(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	h.Add(mkHeapObj("foo", 10)) // nolint: errcheck
	h.Add(mkHeapObj("bar", 1))  // nolint: errcheck
	h.Add(mkHeapObj("bal", 31)) // nolint: errcheck
	h.Add(mkHeapObj("baz", 11)) // nolint: errcheck

	obj, exists, err := h.GetByKey("baz")
	if err != nil || exists == false || obj.(testHeapObject).val != 11 {
		t.Fatalf("unexpected error in getting element")
	}
	// Get non-existing object.
	_, exists, err = h.GetByKey("non-existing")
	if err != nil || exists == true {
		t.Fatalf("didn't expect to get any object")
	}
}

// TestHeap_Close tests Heap.Close and Heap.IsClosed functions.
func TestHeap_Close(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	h.Add(mkHeapObj("foo", 10)) // nolint: errcheck
	h.Add(mkHeapObj("bar", 1))  // nolint: errcheck

	if h.IsClosed() {
		t.Fatalf("didn't expect heap to be closed")
	}
	h.Close()
	if !h.IsClosed() {
		t.Fatalf("expect heap to be closed")
	}
}

// TestHeap_List tests Heap.List function.
func TestHeap_List(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	list := h.List()
	if len(list) != 0 {
		t.Errorf("expected an empty list")
	}

	items := map[string]int{
		"foo": 10,
		"bar": 1,
		"bal": 30,
		"baz": 11,
		"faz": 30,
	}
	for k, v := range items {
		h.Add(mkHeapObj(k, v)) // nolint: errcheck
	}
	list = h.List()
	if len(list) != len(items) {
		t.Errorf("expected %d items, got %d", len(items), len(list))
	}
	for _, obj := range list {
		heapObj := obj.(testHeapObject)
		v, ok := items[heapObj.name]
		if !ok || v != heapObj.val {
			t.Errorf("unexpected item in the list: %v", heapObj)
		}
	}
}

// TestHeap_ListKeys tests Heap.ListKeys function. Scenario is the same as
// TestHeap_list.
func TestHeap_ListKeys(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	list := h.ListKeys()
	if len(list) != 0 {
		t.Errorf("expected an empty list")
	}

	items := map[string]int{
		"foo": 10,
		"bar": 1,
		"bal": 30,
		"baz": 11,
		"faz": 30,
	}
	for k, v := range items {
		h.Add(mkHeapObj(k, v)) // nolint: errcheck
	}
	list = h.ListKeys()
	if len(list) != len(items) {
		t.Errorf("expected %d items, got %d", len(items), len(list))
	}
	for _, key := range list {
		_, ok := items[key]
		if !ok {
			t.Errorf("unexpected item in the list: %v", key)
		}
	}
}

// TestHeapAddAfterClose tests that heap returns an error if anything is added
// after it is closed.
func TestHeapAddAfterClose(t *testing.T) {
	h := New(testHeapObjectKeyFunc, compareInts)
	h.Close()
	if err := h.Add(mkHeapObj("test", 1)); err == nil || err != ErrHeapClosed {
		t.Errorf("expected heap closed error")
	}
	if err := h.AddIfNotPresent(mkHeapObj("test", 1)); err == nil || err != ErrHeapClosed {
		t.Errorf("expected heap closed error")
	}
	if err := h.BulkAdd([]interface{}{mkHeapObj("test", 1)}); err == nil || err != ErrHeapClosed {
		t.Errorf("expected heap closed error")
	}
}
