package fifo

import (
	"errors"
	"sync"

	"github.com/things-go/container"
)

// ErrFIFOClosed used when FIFO is closed.
var ErrFIFOClosed = errors.New("fifo: manipulating with closed queue")

// PopProcessFunc is passed to Pop() method of Queue interface.
// It is supposed to process the accumulator popped from the queue.
type PopProcessFunc[T any] func(T) error

// ErrRequeue may be returned by a PopProcessFunc to safely requeue
// the current item. The value of Err will be returned from Pop.
type ErrRequeue struct {
	// Err is returned by the Pop function
	Err error
}

// Error implement error interface
func (e ErrRequeue) Error() string {
	if e.Err == nil {
		return "the popped item should be requeue without returning an error"
	}
	return e.Err.Error()
}

// Queue extends Store with a collection of Store keys to "process".
// Every Push, Update, or Delete may put the object's key in that collection.
// A Queue has a way to derive the corresponding key given an accumulator.
// A Queue can be accessed concurrently from multiple goroutines.
// A Queue can be "closed", after which Pop operations return an error.
type Queue[T any] interface {
	container.Store[T]

	// Pop blocks until there is at least one key to process or the
	// Queue is closed.  In the latter case Pop returns with an error.
	// In the former case Pop atomically picks one key to process,
	// removes that (key, accumulator) association from the Store, and
	// processes the accumulator.  Pop returns the accumulator that
	// was processed and the result of processing.  The PopProcessFunc
	// may return an ErrRequeue{inner} and in this case Pop will (a)
	// return that (key, accumulator) association to the Queue as part
	// of the atomic processing and (b) return the inner error from
	// Pop.
	Pop(PopProcessFunc[T]) (T, error)

	// AddIfNotPresent puts the given accumulator into the Queue (in
	// association with the accumulator's key) if and only if that key
	// is not already associated with a non-empty accumulator.
	AddIfNotPresent(T) error

	// HasSynced returns true if the first batch of keys have all been
	// popped.  The first batch of keys are those of the first Replace
	// operation if that happened before any Push, AddIfNotPresent,
	// Update, or Delete; otherwise the first batch is empty.
	HasSynced() bool

	// Close the queue
	Close()
}

// Pop is helper function for popping from Queue.
// WARNING: Do NOT use this function in non-test code to avoid races
// unless you really really really really know what you are doing.
func Pop[T any](queue Queue[T]) T {
	var result T
	queue.Pop(func(obj T) error { // nolint: errcheck
		result = obj
		return nil
	})
	return result
}

// FIFO is a Queue
var _ Queue[int] = (*FIFO[int])(nil)

// FIFO is a Queue in which (a) each accumulator is simply the most
// recently provided object and (b) the collection of keys to process
// is a FIFO. The accumulators all start out empty, and deleting an
// object from its accumulator empties the accumulator.  The Resync
// operation is a no-op.
//
// Thus: if multiple adds/updates of a single object happen while that
// object's key is in the queue before it has been processed then it
// will only be processed once, and when it is processed the most
// recent version will be processed. This can't be done with a channel
//
// FIFO solves this use case:
//   - You want to process every object (exactly) once.
//   - You want to process the most recent version of the object when you process it.
//   - You do not want to process deleted objects, they should be removed from the queue.
//   - You do not want to periodically reprocess objects.
//
// Compare with DeltaFIFO for other use cases.
type FIFO[T any] struct {
	rw   sync.RWMutex
	cond sync.Cond
	// We depend on the property that every key in `items` is also in `queue`
	items map[string]T
	queue []string

	// populated is true if the first batch of items inserted by Replace() has been populated
	// or Delete/Push/Update was called first.
	populated bool
	// initialPopulationCount is the number of items inserted by the first call of Replace()
	initialPopulationCount int

	// keyFunc is used to make the key used for queued item insertion and retrieval, and
	// should be deterministic.
	keyFunc container.KeyFunc[T]

	// Indication the queue is closed.
	// Used to indicate a queue is closed so a control loop can exit when a queue is empty.
	// Currently, not used to gate any of CRED operations.
	closed bool
}

// New returns a Store which can be used to queue up items to process.
// keyFunc is used to make the key used for queued item insertion and retrieval, and should be deterministic.
func New[T any](keyFunc container.KeyFunc[T]) *FIFO[T] {
	f := &FIFO[T]{
		items:   map[string]T{},
		queue:   []string{},
		keyFunc: keyFunc,
	}
	f.cond.L = &f.rw
	return f
}

// Close the queue.
func (f *FIFO[T]) Close() {
	f.rw.Lock()
	defer f.rw.Unlock()
	f.closed = true
	f.cond.Broadcast()
}

// HasSynced returns true if an Push/Update/Delete/AddIfNotPresent are called first,
// or the first batch of items inserted by Replace() has been popped.
func (f *FIFO[T]) HasSynced() bool {
	f.rw.Lock()
	defer f.rw.Unlock()
	return f.populated && f.initialPopulationCount == 0
}

// Add inserts an item, and puts it in the queue.
// The item is only enqueued if it doesn't already exist in the set.
func (f *FIFO[T]) Add(obj T) error {
	key, err := f.keyFunc(obj)
	if err != nil {
		return container.KeyError[T]{Obj: obj, Err: err}
	}
	f.rw.Lock()
	defer f.rw.Unlock()
	f.populated = true
	if _, exists := f.items[key]; !exists {
		f.queue = append(f.queue, key)
	}
	f.items[key] = obj
	f.cond.Broadcast()
	return nil
}

// AddIfNotPresent inserts an item, and puts it in the queue. If the item is already
// present in the set, it is neither enqueued nor added to the set.
//
// This is useful in a single producer/consumer scenario so that the consumer can
// safely retry items without contending with the producer and potentially enqueueing
// stale items.
func (f *FIFO[T]) AddIfNotPresent(obj T) error {
	key, err := f.keyFunc(obj)
	if err != nil {
		return container.KeyError[T]{Obj: obj, Err: err}
	}
	f.rw.Lock()
	defer f.rw.Unlock()
	f.addIfNotPresent(key, obj)
	return nil
}

// addIfNotPresent assumes the fifo lock is already held and adds the provided
// item to the queue under id if it does not already exist.
func (f *FIFO[T]) addIfNotPresent(key string, obj T) {
	f.populated = true
	if _, exists := f.items[key]; exists {
		return
	}

	f.queue = append(f.queue, key)
	f.items[key] = obj
	f.cond.Broadcast()
}

// Update is the same as Add in this implementation.
func (f *FIFO[T]) Update(obj T) error {
	return f.Add(obj)
}

// Delete removes an item. It doesn't add it to the queue, because
// this implementation assumes the consumer only cares about the objects,
// not the order in which they were created/added.
func (f *FIFO[T]) Delete(obj T) error {
	id, err := f.keyFunc(obj)
	if err != nil {
		return container.KeyError[T]{Obj: obj, Err: err}
	}
	f.rw.Lock()
	defer f.rw.Unlock()
	f.populated = true
	delete(f.items, id)
	return err
}

// List returns a list of all the items.
func (f *FIFO[T]) List() []T {
	f.rw.RLock()
	defer f.rw.RUnlock()
	return mapValues(f.items)
}

// ListKeys returns a list of all the keys of the objects currently
// in the FIFO.
func (f *FIFO[T]) ListKeys() []string {
	f.rw.RLock()
	defer f.rw.RUnlock()
	return mapKeys(f.items)
}

// Get returns the requested item, or sets exists=false.
func (f *FIFO[T]) Get(obj T) (item T, exists bool, err error) {
	key, err := f.keyFunc(obj)
	if err != nil {
		return item, false, container.KeyError[T]{Obj: obj, Err: err}
	}
	return f.GetByKey(key)
}

// GetByKey returns the requested item, or sets exists=false.
func (f *FIFO[T]) GetByKey(key string) (item T, exists bool, err error) {
	f.rw.RLock()
	defer f.rw.RUnlock()
	item, exists = f.items[key]
	return item, exists, nil
}

// IsClosed checks if the queue is closed.
func (f *FIFO[T]) IsClosed() bool {
	f.rw.RLock()
	defer f.rw.RUnlock()
	return f.closed
}

// Pop waits until an item is ready and processes it. If multiple items are
// ready, they are returned in the order in which they were added/updated.
// The item is removed from the queue (and the store) before it is processed,
// so if you don't successfully process it, it should be added back with
// AddIfNotPresent(). process function is called under lock, so it is safe
// update data structures in it that need to be in sync with the queue.
func (f *FIFO[T]) Pop(process PopProcessFunc[T]) (T, error) {
	var placeholder T

	f.rw.Lock()
	defer f.rw.Unlock()
	for {
		for len(f.queue) == 0 {
			// When the queue is empty, invocation of Pop() is blocked until new item is enqueued.
			// When Close() is called, the f.closed is set and the condition is broadcasted.
			// Which causes this loop to continue and return from the Pop().
			if f.closed {
				return placeholder, ErrFIFOClosed
			}

			f.cond.Wait()
		}
		key := f.queue[0]
		f.queue = f.queue[1:]
		if f.initialPopulationCount > 0 {
			f.initialPopulationCount--
		}
		item, ok := f.items[key]
		if !ok { // Item may have been deleted subsequently.
			continue
		}
		delete(f.items, key)

		var err error
		if process != nil {
			err = process(item)
			if e, ok := err.(ErrRequeue); ok {
				f.addIfNotPresent(key, item)
				err = e.Err
			}
		}
		return item, err
	}
}

// Replace will delete the contents of 'f', using instead the given map.
// 'f' takes ownership of the map, you should not reference the map again
// after calling this function. f's queue is reset, too; upon return, it
// will contain the items in the map, in no particular order.
func (f *FIFO[T]) Replace(list []T, _ string) error {
	items := make(map[string]T, len(list))
	for _, item := range list {
		key, err := f.keyFunc(item)
		if err != nil {
			return container.KeyError[T]{Obj: item, Err: err}
		}
		items[key] = item
	}

	f.rw.Lock()
	defer f.rw.Unlock()

	if !f.populated {
		f.populated = true
		f.initialPopulationCount = len(items)
	}

	f.items = items
	f.queue = f.queue[:0]
	for key := range items {
		f.queue = append(f.queue, key)
	}
	if len(f.queue) > 0 {
		f.cond.Broadcast()
	}
	return nil
}

// Resync will ensure that every object in the Store has its key in the queue.
// This should be a no-op, because that property is maintained by all operations.
func (f *FIFO[T]) Resync() error {
	f.rw.Lock()
	defer f.rw.Unlock()

	inQueue := make(map[string]struct{}, len(f.queue))
	for _, v := range f.queue {
		inQueue[v] = struct{}{}
	}
	for key := range f.items {
		if _, ok := inQueue[key]; !ok {
			f.queue = append(f.queue, key)
		}
	}
	if len(f.queue) > 0 {
		f.cond.Broadcast()
	}
	return nil
}

// MapValues returns the values of the map m.
// The values will be in an indeterminate order.
func mapValues[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

// mapKeys returns the keys of the map m.
// The keys will be in an indeterminate order.
func mapKeys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}
