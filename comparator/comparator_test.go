// Copyright [2020] [thinkgos]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package comparator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	// nil
	assert.Equal(t, Compare(nil, nil), 0)

	// bool
	assert.Equal(t, Compare(false, true), -1)
	assert.Equal(t, Compare(false, false), 0)
	assert.Equal(t, Compare(true, true), 0)
	assert.Equal(t, Compare(true, false), 1)

	// int
	assert.Equal(t, Compare(1, 2), -1)
	assert.Equal(t, Compare(2, 2), 0)
	assert.Equal(t, Compare(2, 1), 1)

	// int8
	assert.Equal(t, Compare(int8(1), int8(2)), -1)
	assert.Equal(t, Compare(int8(2), int8(2)), 0)
	assert.Equal(t, Compare(int8(2), int8(1)), 1)

	// int16
	assert.Equal(t, Compare(int16(1), int16(2)), -1)
	assert.Equal(t, Compare(int16(2), int16(2)), 0)
	assert.Equal(t, Compare(int16(2), int16(1)), 1)

	// int32
	assert.Equal(t, Compare(int32(1), int32(2)), -1)
	assert.Equal(t, Compare(int32(2), int32(2)), 0)
	assert.Equal(t, Compare(int32(2), int32(1)), 1)

	// rune
	assert.Equal(t, Compare(rune(1), rune(2)), -1)
	assert.Equal(t, Compare(rune(2), rune(2)), 0)
	assert.Equal(t, Compare(rune(2), rune(1)), 1)

	// int64
	assert.Equal(t, Compare(int64(1), int64(2)), -1)
	assert.Equal(t, Compare(int64(2), int64(2)), 0)
	assert.Equal(t, Compare(int64(2), int64(1)), 1)

	// uint
	assert.Equal(t, Compare(uint(1), uint(2)), -1)
	assert.Equal(t, Compare(uint(2), uint(2)), 0)
	assert.Equal(t, Compare(uint(2), uint(1)), 1)

	// uint8
	assert.Equal(t, Compare(uint8(1), uint8(2)), -1)
	assert.Equal(t, Compare(uint8(2), uint8(2)), 0)
	assert.Equal(t, Compare(uint8(2), uint8(1)), 1)

	// byte
	assert.Equal(t, Compare(byte(1), byte(2)), -1)
	assert.Equal(t, Compare(byte(2), byte(2)), 0)
	assert.Equal(t, Compare(byte(2), byte(1)), 1)

	// uint16
	assert.Equal(t, Compare(uint16(1), uint16(2)), -1)
	assert.Equal(t, Compare(uint16(2), uint16(2)), 0)
	assert.Equal(t, Compare(uint16(2), uint16(1)), 1)

	// uint32
	assert.Equal(t, Compare(uint32(1), uint32(2)), -1)
	assert.Equal(t, Compare(uint32(2), uint32(2)), 0)
	assert.Equal(t, Compare(uint32(2), uint32(1)), 1)
	// uint64
	assert.Equal(t, Compare(uint64(1), uint64(2)), -1)
	assert.Equal(t, Compare(uint64(2), uint64(2)), 0)
	assert.Equal(t, Compare(uint64(2), uint64(1)), 1)

	// float32
	assert.Equal(t, Compare(float32(1), float32(2)), -1)
	assert.Equal(t, Compare(float32(2), float32(2)), 0)
	assert.Equal(t, Compare(float32(2), float32(1)), 1)

	// float64
	assert.Equal(t, Compare(float64(1), float64(2)), -1)
	assert.Equal(t, Compare(float64(2), float64(2)), 0)
	assert.Equal(t, Compare(float64(2), float64(1)), 1)

	// string
	assert.Equal(t, Compare("abc", "ade"), -1)
	assert.Equal(t, Compare("ade", "ade"), 0)
	assert.Equal(t, Compare("ade", "abc"), 1)

	// time.Time
	t1, t2 := time.Now(), time.Now().Add(10*time.Hour)
	assert.Equal(t, Compare(t1, t2), -1)
	assert.Equal(t, Compare(t1, t1), 0)
	assert.Equal(t, Compare(t2, t1), 1)

	// cause panic
	assert.Panics(t, func() { Compare(nil, 1) })
	assert.Panics(t, func() { Compare(uint(1), 1) })
	assert.Panics(t, func() { Compare(uint(1), 1) })
	assert.Panics(t, func() { Compare(time.Now(), struct{}{}) })
	assert.Panics(t, func() { Compare(map[string]string{"a": "b"}, map[string]string{"a": "b"}) })
}
