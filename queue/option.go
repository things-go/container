package queue

import (
	"github.com/things-go/container/comparator"
)

// Apply option apply interface.
type Apply interface {
	apply(c comparator.Comparator)
}

// Option option for New.
type Option func(a Apply)

// WithComparator with user's Comparator.
func WithComparator(c comparator.Comparator) Option {
	return func(a Apply) {
		a.apply(c)
	}
}
