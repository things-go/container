// This file implements a heap data structure.
package heap

import (
	"container/heap"
	"errors"
	"fmt"
	"sync"

	"github.com/things-go/container"
	"github.com/things-go/container/comparator"
)

// ErrHeapClosed used when heap is closed.
var ErrHeapClosed = errors.New("heap is closed")

type heapItem struct {
	obj   interface{} // The object which is stored in the heap.
	index int         // The index of the object's key in the Heap.queue.
}

type itemKeyValue struct {
	key string
	obj interface{}
}

// heapData is an internal struct that implements the standard heap interface
// and keeps the data stored in the heap.
type heapData struct {
	// items is a map from key of the objects to the objects and their index.
	// We depend on the property that items in the map are in the queue and vice versa.
	items map[string]*heapItem
	// queue implements a heap data structure and keeps the order of elements
	// according to the heap invariant. The queue keeps the keys of objects stored
	// in "items".
	queue []string

	// compare is used to compare two objects in the heap.
	compare container.Comparator
}

var _ heap.Interface = (*heapData)(nil) // heapData is a standard heap

// Less compares two objects and returns true if the first one should go
// in front of the second one in the heap.
func (h *heapData) Less(i, j int) bool {
	if i > len(h.queue) || j > len(h.queue) {
		return false
	}
	itemi, ok := h.items[h.queue[i]]
	if !ok {
		return false
	}
	itemj, ok := h.items[h.queue[j]]
	if !ok {
		return false
	}
	return h.Compare(itemi.obj, itemj.obj) < 0
}

// Len returns the number of items in the Heap.
func (h *heapData) Len() int { return len(h.queue) }

// Swap implements swapping of two elements in the heap. This is a part of standard
// heap interface and should never be called directly.
func (h *heapData) Swap(i, j int) {
	h.queue[i], h.queue[j] = h.queue[j], h.queue[i]
	item := h.items[h.queue[i]]
	item.index = i
	item = h.items[h.queue[j]]
	item.index = j
}

// Push is supposed to be called by heap.Push only.
func (h *heapData) Push(kv interface{}) {
	keyValue := kv.(*itemKeyValue)
	n := len(h.queue)
	h.items[keyValue.key] = &heapItem{keyValue.obj, n}
	h.queue = append(h.queue, keyValue.key)
}

// Pop is supposed to be called by heap.Pop only.
func (h *heapData) Pop() interface{} {
	key := h.queue[len(h.queue)-1]
	h.queue = h.queue[0 : len(h.queue)-1]
	item, ok := h.items[key]
	if !ok {
		// This is an error
		return nil
	}
	delete(h.items, key)
	return item.obj
}

func (h *heapData) Compare(v1, v2 interface{}) int {
	if h.compare == nil {
		return comparator.Compare(v1, v2)
	}
	return h.compare(v1, v2)
}

// Contain return the requested item exist or not with key
func (h *heapData) Contain(key string) bool {
	_, exist := h.items[key]
	return exist
}

// GetItem return the requested item with key
func (h *heapData) GetItem(key string) (*heapItem, bool) {
	v, exist := h.items[key]
	return v, exist
}

// ListKeys returns a list of all the keys of the objects currently in the heapData.
func (h *heapData) ListKeys() []string {
	list := make([]string, 0, len(h.items))
	for key := range h.items {
		list = append(list, key)
	}
	return list
}

// List returns a list of all the items obj.
func (h *heapData) List() []interface{} {
	list := make([]interface{}, 0, len(h.items))
	for _, item := range h.items {
		list = append(list, item.obj)
	}
	return list
}

// Heap is a thread-safe producer/consumer queue that implements a heap data structure.
// It can be used to implement priority queues and similar data structures.
type Heap struct {
	rw   sync.RWMutex
	cond sync.Cond

	// data stores objects and has a queue that keeps their ordering according
	// to the heap invariant.
	data *heapData

	// keyFunc is used to make the key used for queued item insertion and retrieval, and
	// should be deterministic.
	keyFunc container.KeyFunc

	// closed indicates that the queue is closed.
	// It is mainly used to let Pop() exit its control loop while waiting for an item.
	closed bool
}

// New returns a Heap which can be used to queue up items to process.
func New(keyFn container.KeyFunc, cmp container.Comparator) *Heap {
	h := &Heap{
		data: &heapData{
			items:   map[string]*heapItem{},
			queue:   []string{},
			compare: cmp,
		},
		keyFunc: keyFn,
	}
	h.cond.L = &h.rw
	return h
}

// Close the Heap and signals condition variables that may be waiting to pop
// items from the heap.
func (sf *Heap) Close() {
	sf.rw.Lock()
	defer sf.rw.Unlock()
	sf.closed = true
	sf.cond.Broadcast()
}

// IsClosed returns true if the queue is closed.
func (sf *Heap) IsClosed() bool {
	sf.rw.RLock()
	defer sf.rw.RUnlock()
	return sf.closed
}

// Add inserts an item, and puts it in the queue. The item is updated if it
// already exists.
func (sf *Heap) Add(obj interface{}) error {
	key, err := sf.keyFunc(obj)
	if err != nil {
		return container.KeyError{Obj: obj, Err: err}
	}
	sf.rw.Lock()
	defer sf.rw.Unlock()
	if sf.closed {
		return ErrHeapClosed
	}
	if vv, exists := sf.data.GetItem(key); exists {
		vv.obj = obj
		heap.Fix(sf.data, vv.index)
	} else {
		sf.addIfNotPresentLocked(key, obj)
	}
	sf.cond.Broadcast()
	return nil
}

// BulkAdd adds all the items in the list to the queue and then signals the condition
// variable. It is useful when the caller would like to add all of the items
// to the queue before consumer starts processing them.
func (sf *Heap) BulkAdd(list []interface{}) error {
	sf.rw.Lock()
	defer sf.rw.Unlock()
	if sf.closed {
		return ErrHeapClosed
	}
	for _, obj := range list {
		key, err := sf.keyFunc(obj)
		if err != nil {
			return container.KeyError{Obj: obj, Err: err}
		}

		if vv, exists := sf.data.GetItem(key); exists {
			vv.obj = obj
			heap.Fix(sf.data, vv.index)
		} else {
			sf.addIfNotPresentLocked(key, obj)
		}
	}
	sf.cond.Broadcast()
	return nil
}

// AddIfNotPresent inserts an item, and puts it in the queue. If an item with
// the key is present in the map, no changes is made to the item.
//
// This is useful in a single producer/consumer scenario so that the consumer can
// safely retry items without contending with the producer and potentially enqueueing
// stale items.
func (sf *Heap) AddIfNotPresent(obj interface{}) error {
	key, err := sf.keyFunc(obj)
	if err != nil {
		return container.KeyError{Obj: obj, Err: err}
	}
	sf.rw.Lock()
	defer sf.rw.Unlock()
	if sf.closed {
		return ErrHeapClosed
	}
	sf.addIfNotPresentLocked(key, obj)
	sf.cond.Broadcast()
	return nil
}

// addIfNotPresentLocked assumes the lock is already held and adds the provided
// item to the queue if it does not already exist.
func (sf *Heap) addIfNotPresentLocked(key string, obj interface{}) {
	if sf.data.Contain(key) {
		return
	}
	heap.Push(sf.data, &itemKeyValue{key, obj})
}

// Update is the same as Push in this implementation. When the item does not
// exist, it is added.
func (sf *Heap) Update(obj interface{}) error {
	return sf.Add(obj)
}

// Delete removes an item.
func (sf *Heap) Delete(obj interface{}) error {
	key, err := sf.keyFunc(obj)
	if err != nil {
		return container.KeyError{Obj: obj, Err: err}
	}
	sf.rw.Lock()
	defer sf.rw.Unlock()
	if item, ok := sf.data.GetItem(key); ok {
		heap.Remove(sf.data, item.index)
		return nil
	}
	return errors.New("object not found")
}

// Pop waits until an item is ready. If multiple items are
// ready, they are returned in the order given by Heap.data.compare.
func (sf *Heap) Pop() (interface{}, error) {
	sf.rw.Lock()
	defer sf.rw.Unlock()
	for len(sf.data.queue) == 0 {
		// When the queue is empty, invocation of Pop() is blocked until new item is enqueued.
		// When Close() is called, the h.closed is set and the condition is broadcast,
		// which causes this loop to continue and return from the Pop().
		if sf.closed {
			return nil, ErrHeapClosed
		}
		sf.cond.Wait()
	}
	obj := heap.Pop(sf.data)
	if obj == nil {
		return nil, fmt.Errorf("object was removed from heap data")
	}

	return obj, nil
}

// List returns a list of all the items.
func (sf *Heap) List() []interface{} {
	sf.rw.RLock()
	defer sf.rw.RUnlock()
	return sf.data.List()
}

// ListKeys returns a list of all the keys of the objects currently in the Heap.
func (sf *Heap) ListKeys() []string {
	sf.rw.RLock()
	defer sf.rw.RUnlock()
	return sf.data.ListKeys()
}

// Get returns the requested item, or sets exists=false.
func (sf *Heap) Get(obj interface{}) (interface{}, bool, error) {
	key, err := sf.keyFunc(obj)
	if err != nil {
		return nil, false, container.KeyError{Obj: obj, Err: err}
	}
	return sf.GetByKey(key)
}

// GetByKey returns the requested item, or sets exists=false.
func (sf *Heap) GetByKey(key string) (interface{}, bool, error) {
	sf.rw.RLock()
	defer sf.rw.RUnlock()
	item, exists := sf.data.GetItem(key)
	if !exists {
		return nil, false, nil
	}
	return item.obj, true, nil
}
