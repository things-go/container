package container

// Comparator compares its two arguments for order.
// It returns a negative integer, zero,
// or a positive integer as the first argument
// is less than, equal to, or greater than the second.
type Comparator func(v1, v2 interface{}) int
