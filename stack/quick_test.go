package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuickStack(t *testing.T) {
	s := NewQuickStack[string]()
	s.Push("5")
	s.Push("hello")

	// length
	assert.Equal(t, 2, s.Len())
	assert.False(t, s.IsEmpty())

	// Peek "hello"
	val1, ok := s.Peek()
	assert.True(t, ok)
	assert.Equal(t, "hello", val1)

	// Pop "hello"
	val2, ok := s.Pop()
	assert.True(t, ok)
	assert.Equal(t, "hello", val2)

	// Peek 5
	val3, ok := s.Peek()
	assert.True(t, ok)
	assert.Equal(t, "5", val3)

	// Pop 5
	val4, ok := s.Pop()
	assert.True(t, ok)
	assert.Equal(t, "5", val4)

	val5, ok := s.Pop()
	assert.False(t, ok)
	assert.Zero(t, val5)

	val6, ok := s.Peek()
	assert.False(t, ok)
	assert.Zero(t, val6)

	s.Push("5")
	s.Push("6")

	s1 := s.Clone()
	assert.Equal(t, 2, s1.Len())

	assert.False(t, s.IsEmpty())
	s.Clear()
	assert.True(t, s.IsEmpty())
	assert.Zero(t, s.Len())
}
