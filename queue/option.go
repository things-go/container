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

// WithComparator compare use custom Comparator.
// if not set, use reflect.DeepEqual
func WithComparator(c container.Comparator) Option {
	return func(a Apply) {
		a.apply(c)
	}
}
