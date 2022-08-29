package container

import (
	"fmt"
)

// Store is a generic object storage and processing interface.  A
// Store holds a map from string keys to accumulators, and has
// operations to add, update, and delete a given object to/from the
// accumulator currently associated with a given key.  A Store also
// knows how to extract the key from a given object, so many operations
// are given only the object.
//
// In the simplest Store implementations each accumulator is simply
// the last given object, or empty after Delete, and thus the Store's
// behavior is simple storage.
//
// Reflector knows how to watch a server and update a Store.  This
// package provides a variety of implementations of Store.
type Store[T any] interface {

	// Add adds the given object to the accumulator associated with the given object's key
	Add(obj T) error

	// Update updates the given object in the accumulator associated with the given object's key
	Update(obj T) error

	// Delete deletes the given object from the accumulator associated with the given object's key
	Delete(obj T) error

	// List returns a list of all the currently non-empty accumulators
	List() []T

	// ListKeys returns a list of all the keys currently associated with non-empty accumulators
	ListKeys() []string

	// Get returns the accumulator associated with the given object's key
	Get(obj T) (item T, exists bool, err error)

	// GetByKey returns the accumulator associated with the given key
	GetByKey(key string) (item T, exists bool, err error)

	// Replace will delete the contents of the store, using instead the
	// given list. Store takes ownership of the list, you should not reference
	// it after calling this function.
	Replace([]T, string) error

	// Resync is meaningless in the terms appearing here but has
	// meaning in some implementations that have non-trivial
	// additional behavior (e.g., DeltaFIFO).
	Resync() error
}

// KeyFunc knows how to make a key from an object. Implementations should be deterministic.
type KeyFunc[T any] func(obj T) (string, error)

// KeyError will be returned any time a KeyFunc gives an error; it includes the object
// at fault.
type KeyError[T any] struct {
	Obj T
	Err error
}

// Error gives a human-readable description of the error.
func (k KeyError[T]) Error() string {
	return fmt.Sprintf("couldn't create key for object %+v: %v", k.Obj, k.Err)
}
