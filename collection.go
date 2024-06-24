package container

type Comparable interface {
	CompareTo(Comparable) int
}

// Stack is a Stack interface, which is LIFO (last-in-first-out).
type Stack[T any] interface {
	// Len returns the number of elements in the collection.
	Len() int
	// IsEmpty returns true if this container contains no elements.
	IsEmpty() bool
	// Clear initializes or clears all of the elements from this container.
	Clear()
	// Push pushes an element into this Stack.
	Push(T)
	// Pop pops the element on the top of this Stack.
	Pop() (T, bool)
	// Peek retrieves, but does not remove, the element on the top of this Stack, or return false if this Stack is empty.
	Peek() (T, bool)
}

// Queue is a type of Queue, which is FIFO(first-in-first-out).
type Queue[T comparable] interface {
	// Len returns the number of elements in the collection.
	Len() int
	// IsEmpty returns true if this container contains no elements.
	IsEmpty() bool
	// Clear initializes or clears all of the elements from this container.
	Clear()
	// Add inserts an element into the tail of this Queue.
	Add(T)
	// Peek retrieves, but does not remove, the head of this Queue, or return nil if this Queue is empty.
	Peek() (T, bool)
	// Poll retrieves and removes the head of the this Queue, or return nil if this Queue is empty.
	Poll() (T, bool)
	// Remove a single instance of the specified element from this queue, if it is present.
	Remove(val T)
	// Contains returns true if this queue contains the specified element.
	Contains(val T) bool
}

// List is a type of list, both ArrayList and LinkedList implement this interface.
type List[T comparable] interface {
	// Len returns the number of elements in the collection.
	Len() int
	// IsEmpty returns true if this container contains no elements.
	IsEmpty() bool
	// Clear initializes or clears all of the elements from this container.
	Clear()
	// Push appends the specified elements to the end of this list.
	Push(vals T)
	// PushFront inserts a new element e with value v at the front of list l
	PushFront(v T)
	// PushBack inserts a new element e with value v at the back of list l.
	PushBack(v T)
	// Add inserts the specified element at the specified position in this list.
	Add(index int, val T) error

	// Poll return the front element value and then remove from list
	Poll() (T, bool)
	// PollFront return the front element value and then remove from list
	PollFront() (T, bool)
	// PollBack return the back element value and then remove from list
	PollBack() (T, bool)
	// Remove removes the element at the specified position in this list.
	// It returns an error if the index is out of range.
	Remove(index int) (T, error)
	// RemoveValue removes the first occurrence of the specified element from this list, if it is present.
	// It returns false if the target value isn't present, otherwise returns true.
	RemoveValue(val T) bool

	// Get returns the element at the specified position in this list. The index must be in the range of [0, size).
	Get(index int) (T, error)
	// Peek return the front element value
	Peek() (T, bool)
	// PeekFront return the front element value
	PeekFront() (T, bool)
	// PeekBack return the back element value
	PeekBack() (T, bool)

	// Iterator returns an iterator over the elements in this list in proper sequence.
	Iterator(f func(T) bool)
	// ReverseIterator returns an iterator over the elements in this list in reverse sequence as Iterator.
	ReverseIterator(f func(T) bool)

	// Contains returns true if this list contains the specified element.
	Contains(val T) bool
	// Sort sorts the element using default options below.
	// It sorts the elements into ascending sequence according to their natural ordering.
	Sort(less func(a, b T) int)
	// Values get a copy of all the values in the list
	Values() []T
}

// LinkedMap is a type of linked map, and LinkedMap implements this interface.
type LinkedMap[K comparable, V any] interface {
	// Cap returns the capacity of elements of list l.
	Cap() int
	// Len returns the number of elements in the collection.
	Len() int
	// IsEmpty returns true if this container contains no elements.
	IsEmpty() bool
	// Clear initializes or clears all of the elements from this container.
	Clear()
	// Push associates the specified value with the specified key in this map.
	// If the map previously contained a mapping for the key,
	// the old value is replaced by the specified value. and then move the item to the back of the list.
	// If over the cap, it will remove the back item then push new item to back
	// It returns the previous value associated with the specified key, or nil if there was no mapping for the key.
	// A nil return can also indicate that the map previously associated nil with the specified key.
	Push(k K, v V) (V, bool)
	// PushFront associates the specified value with the specified key in this map.
	// If the map previously contained a mapping for the key,
	// the old value is replaced by the specified value. and then move the item to the front of the list.
	// If over the cap, it will remove the back item then push new item to front
	// It returns the previous value associated with the specified key, or nil if there was no mapping for the key.
	// A nil return can also indicate that the map previously associated nil with the specified key.
	PushFront(k K, v V) (V, bool)
	// PushBack associates the specified value with the specified key in this map.
	// If the map previously contained a mapping for the key,
	// the old value is replaced by the specified value. and then move the item to the back of the list.
	// If over the cap, it will remove the back item then push new item to back
	PushBack(k K, v V) (V, bool)

	// Poll removes the first element from this map, which is the head of the list.
	// It returns the (key, value, true) if the map isn't empty, or (nil, nil, false) if the map is empty.
	Poll() (K, V, bool)
	// PollFront return the front element value and then remove from list
	PollFront() (k K, v V, exist bool)
	// PollBack removes the last element from this map, which is the tail of the list.
	// It returns the (key, value, true) if the map isn't empty, or (nil, nil, false) if the map is empty.
	PollBack() (K, V, bool)
	// Remove removes the mapping for a key from this map if it is present.
	// It returns the value to which this map previously associated the key, and true,
	// or nil and false if the map contained no mapping for the key.
	Remove(k K) (V, bool)

	// Get returns the value to which the specified key is mapped, or nil if this map contains no mapping for the key.
	Get(k K, defaultValue ...V) V
	// Peek return the front element value
	Peek() (k K, v V, exist bool)
	// PeekFront return the front element value
	PeekFront() (k K, v V, exist bool)
	// PeekBack return the back element value
	PeekBack() (k K, v V, exist bool)

	// Iterator returns an iterator over the elements in this map in proper sequence.
	Iterator(cb func(k K, v V) bool)
	// ReverseIterator returns an iterator over the elements in this map in reverse sequence as Iterator.
	ReverseIterator(cb func(k K, v V) bool)

	// Contains returns true if this map contains a mapping for the specified key.
	Contains(k K) bool
	// ContainsValue returns true if this map maps one or more keys to the specified value.
	ContainsValue(v V, equal func(a, b V) bool) bool
}
