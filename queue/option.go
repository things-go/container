package queue

import (
	"github.com/things-go/container"
)

// Apply option apply interface.
type Apply interface {
	apply(c container.Comparator)
}

// Option option for New.
type Option func(a Apply)

// WithComparator with user's Comparator.
func WithComparator(c container.Comparator) Option {
	return func(a Apply) {
		a.apply(c)
	}
}
