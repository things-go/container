package stack

import (
	"testing"
)

func BenchmarkStack(b *testing.B) {
	q := New()
	for i := 0; i < b.N; i++ {
		q.Push(1)
		q.Pop()
	}
}

func BenchmarkQuickStack(b *testing.B) {
	q := NewQuickStack()
	for i := 0; i < b.N; i++ {
		q.Push(1)
		q.Pop()
	}
}
