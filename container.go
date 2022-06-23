package container

import (
	"golang.org/x/exp/constraints"

	"github.com/things-go/container/core/heap"
)

var _ heap.Interface[int] = (*Container[int])(nil)

type Container[T constraints.Ordered] struct {
	Items   []T
	Reverse bool
}

// Len implement heap.Interface.
func (c Container[T]) Len() int {
	return len(c.Items)
}

// Swap implement heap.Interface.
func (c Container[T]) Swap(i, j int) {
	c.Items[i], c.Items[j] = c.Items[j], c.Items[i]
}

// Less implement heap.Interface.
func (c Container[T]) Less(i, j int) bool {
	if c.Reverse {
		i, j = j, i
	}
	return c.Items[i] < c.Items[j]
}

// Push implement heap.Interface.
func (c *Container[T]) Push(x T) {
	c.Items = append(c.Items, x)
}

// Pop implement heap.Interface.
func (c *Container[T]) Pop() T {
	var placeholder T

	old := c.Items
	n := len(old)
	x := old[n-1]
	old[n-1] = placeholder
	c.Items = old[:n-1]
	return x
}
