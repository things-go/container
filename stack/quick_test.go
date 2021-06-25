package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuickStack(t *testing.T) {
	s := NewQuickStack()
	s.Push(5)
	s.Push("hello")

	// length
	assert.Equal(t, 2, s.Len())
	assert.False(t, s.IsEmpty())

	// Peek "hello"
	val1, ok := s.Peek().(string)
	assert.True(t, ok)
	assert.Equal(t, "hello", val1)

	// Pop "hello"
	val2, ok := s.Pop().(string)
	assert.True(t, ok)
	assert.Equal(t, "hello", val2)

	// Peek 5
	val3, ok := s.Peek().(int)
	assert.True(t, ok)
	assert.Equal(t, 5, val3)

	// Pop 5
	val4, ok := s.Pop().(int)
	assert.True(t, ok)
	assert.Equal(t, 5, val4)

	val5 := s.Pop()
	assert.Nil(t, val5)
	assert.Nil(t, s.Peek())

	s.Push(5)
	s.Push(6)
	assert.False(t, s.IsEmpty())
	s.Clear()
	assert.True(t, s.IsEmpty())
	assert.Zero(t, s.Len())
}
