package queue

import (
	"testing"
)

func BenchmarkQueue(b *testing.B) {
	q := New[int]()
	for i := 0; i < b.N; i++ {
		q.Add(1)
		q.Poll()
	}
}

func BenchmarkQuickQueue(b *testing.B) {
	q := NewQuickQueue[int]()
	for i := 0; i < b.N; i++ {
		q.Add(1)
		q.Poll()
	}
}
