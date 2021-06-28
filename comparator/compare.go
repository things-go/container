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
	"fmt"
	"reflect"
	"time"
)

// Compare compares its two arguments
// if they have the same type and are comparable,
// otherwise returns an error in the second return value.
// It returns a negative integer, zero,
// or a positive integer as the first argument is
// less than, equal to, or greater than the second.
func Compare(v1, v2 interface{}) int {
	if v1 == nil && v2 == nil {
		return 0
	}
	if v1 == nil || v2 == nil {
		panic("A nil value can't be compared to a non-nil value")
	}

	k1, k2 := reflect.TypeOf(v1).Kind(), reflect.TypeOf(v2).Kind()
	if k1 != k2 {
		panic(fmt.Sprintf("Two values of different type can't be compared, %s: %s", k1, k2))
	}

	ret := 0
	switch k1 { // nolint: exhaustive
	case reflect.Int:
		cv1, cv2 := v1.(int), v2.(int)
		if cv1 < cv2 {
			ret = -1
		} else if cv1 > cv2 {
			ret = 1
		}
	case reflect.Uint:
		cv1, cv2 := v1.(uint), v2.(uint)
		if cv1 < cv2 {
			ret = -1
		} else if cv1 > cv2 {
			ret = 1
		}
	case reflect.String:
		cv1, cv2 := v1.(string), v2.(string)
		if cv1 < cv2 {
			ret = -1
		} else if cv1 > cv2 {
			ret = 1
		}
	case reflect.Float32:
		cv1, cv2 := v1.(float32), v2.(float32)
		if cv1 < cv2 {
			ret = -1
		} else if cv1 > cv2 {
			ret = 1
		}
	case reflect.Float64:
		cv1, cv2 := v1.(float64), v2.(float64)
		if cv1 < cv2 {
			ret = -1
		} else if cv1 > cv2 {
			ret = 1
		}
	case reflect.Bool:
		// false < true
		b1, b2 := v1.(bool), v2.(bool)
		if !b1 && b2 { // b1 == false && b2 == true
			ret = -1
		} else if b1 && !b2 { // b1 == true && b2 == false
			ret = 1
		}
	case reflect.Int8:
		cv1, cv2 := v1.(int8), v2.(int8)
		if cv1 < cv2 {
			ret = -1
		} else if cv1 > cv2 {
			ret = 1
		}
	case reflect.Uint8: // valid for both uint8 and byte
		cv1, cv2 := v1.(uint8), v2.(uint8)
		if cv1 < cv2 {
			ret = -1
		} else if cv1 > cv2 {
			ret = 1
		}
	case reflect.Int16:
		cv1, cv2 := v1.(int16), v2.(int16)
		if cv1 < cv2 {
			ret = -1
		} else if cv1 > cv2 {
			ret = 1
		}
	case reflect.Uint16:
		cv1, cv2 := v1.(uint16), v2.(uint16)
		if cv1 < cv2 {
			ret = -1
		} else if cv1 > cv2 {
			ret = 1
		}
	case reflect.Int32: // valid for both int32 and rune
		cv1, cv2 := v1.(int32), v2.(int32)
		if cv1 < cv2 {
			ret = -1
		} else if cv1 > cv2 {
			ret = 1
		}
	case reflect.Uint32:
		cv1, cv2 := v1.(uint32), v2.(uint32)
		if cv1 < cv2 {
			ret = -1
		} else if cv1 > cv2 {
			ret = 1
		}
	case reflect.Int64:
		cv1, cv2 := v1.(int64), v2.(int64)
		if cv1 < cv2 {
			ret = -1
		} else if cv1 > cv2 {
			ret = 1
		}
	case reflect.Uint64:
		cv1, cv2 := v1.(uint64), v2.(uint64)
		if cv1 < cv2 {
			ret = -1
		} else if cv1 > cv2 {
			ret = 1
		}
	case reflect.Struct:
		// compare time
		timeCmpRet, isBothTime := CompareTime(v1, v2)
		if isBothTime {
			return timeCmpRet
		}
		panic("please define a customized comparator.Comparator for your struct")
	default:
		panic(fmt.Sprintf("type '%s' can't be compared,please define a customized comparator.Comparator for your struct", k1))
	}
	return ret
}

// CompareTime compares its two arguments if both of them are time.Time, and returns true
// and the comparison result; otherwise return false in the first return argument.
func CompareTime(v1, v2 interface{}) (int, bool) {
	t1, ok1 := v1.(time.Time)
	t2, ok2 := v2.(time.Time)
	if ok1 && ok2 {
		if t1.Before(t2) {
			return -1, true
		}
		if t1.After(t2) {
			return 1, true
		}
		return 0, true
	}
	return 0, false
}
