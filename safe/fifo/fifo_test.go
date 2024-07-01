package fifo

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func testFifoObjectKeyFunc(obj testFifoObject) (string, error) {
	return obj.name, nil
}

type testFifoObject struct {
	name string
	val  interface{}
}

func mkFifoObj(name string, val interface{}) testFifoObject {
	return testFifoObject{name: name, val: val}
}

func Test_FIFO_basic(t *testing.T) {
	f := New[testFifoObject](testFifoObjectKeyFunc)
	const amount = 500
	go func() {
		for i := 0; i < amount; i++ {
			f.Add(mkFifoObj(string([]rune{'a', rune(i)}), i+1)) // nolint: errcheck
		}
	}()
	go func() {
		for u := uint64(0); u < amount; u++ {
			f.Add(mkFifoObj(string([]rune{'b', rune(u)}), u+1)) // nolint: errcheck
		}
	}()

	lastInt := int(0)
	lastUint := uint64(0)
	for i := 0; i < amount*2; i++ {
		switch obj := Pop[testFifoObject](f).val.(type) {
		case int:
			if obj <= lastInt {
				t.Errorf("got %v (int) out of order, last was %v", obj, lastInt)
			}
			lastInt = obj
		case uint64:
			if obj <= lastUint {
				t.Errorf("got %v (uint) out of order, last was %v", obj, lastUint)
			} else {
				lastUint = obj
			}
		default:
			t.Fatalf("unexpected type %#v", obj)
		}
	}
}

func Test_FIFO_requeueOnPop(t *testing.T) {
	f := New[testFifoObject](testFifoObjectKeyFunc)

	f.Add(mkFifoObj("foo", 10)) // nolint: errcheck
	_, err := f.Pop(func(obj testFifoObject) error {
		if obj.name != "foo" {
			t.Fatalf("unexpected object: %#v", obj)
		}
		return ErrRequeue{Err: nil}
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok, err := f.GetByKey("foo"); err != nil || !ok {
		t.Fatalf("object should have been requeued: %t %v", ok, err)
	}

	_, err = f.Pop(func(obj testFifoObject) error {
		if obj.name != "foo" {
			t.Fatalf("unexpected object: %#v", obj)
		}
		return ErrRequeue{Err: fmt.Errorf("test error")}
	})
	if err == nil || err.Error() != "test error" {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok, err := f.GetByKey("foo"); err != nil || !ok {
		t.Fatalf("object should have been requeued: %t %v", ok, err)
	}

	_, err = f.Pop(func(obj testFifoObject) error {
		if obj.name != "foo" {
			t.Fatalf("unexpected object: %#v", obj)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok, err := f.GetByKey("foo"); ok || err != nil {
		t.Fatalf("object should have been removed: %t %v", ok, err)
	}
}

func Test_FIFO_addUpdate(t *testing.T) {
	f := New[testFifoObject](testFifoObjectKeyFunc)
	f.Add(mkFifoObj("foo", 10))    // nolint: errcheck
	f.Update(mkFifoObj("foo", 15)) // nolint: errcheck

	if e, a := []testFifoObject{mkFifoObj("foo", 15)}, f.List(); !reflect.DeepEqual(e, a) {
		t.Errorf("Expected %+v, got %+v", e, a)
	}
	if e, a := []string{"foo"}, f.ListKeys(); !reflect.DeepEqual(e, a) {
		t.Errorf("Expected %+v, got %+v", e, a)
	}

	got := make(chan testFifoObject, 2)
	go func() {
		for {
			got <- Pop[testFifoObject](f)
		}
	}()

	first := <-got
	if e, a := 15, first.val; e != a {
		t.Errorf("Didn't get updated value (%v), got %v", e, a)
	}
	select {
	case unexpected := <-got:
		t.Errorf("Got second value %v", unexpected.val)
	case <-time.After(50 * time.Millisecond):
	}
	_, exists, _ := f.Get(mkFifoObj("foo", ""))
	if exists {
		t.Errorf("item did not get removed")
	}
}

func Test_FIFO_addReplace(t *testing.T) {
	f := New(testFifoObjectKeyFunc)
	f.Add(mkFifoObj("foo", 10))                             // nolint: errcheck
	f.Replace([]testFifoObject{mkFifoObj("foo", 15)}, "15") // nolint: errcheck
	got := make(chan testFifoObject, 2)
	go func() {
		for {
			got <- Pop[testFifoObject](f)
		}
	}()

	first := <-got
	if e, a := 15, first.val; e != a {
		t.Errorf("Didn't get updated value (%v), got %v", e, a)
	}
	select {
	case unexpected := <-got:
		t.Errorf("Got second value %v", unexpected.val)
	case <-time.After(50 * time.Millisecond):
	}
	_, exists, _ := f.Get(mkFifoObj("foo", ""))
	if exists {
		t.Errorf("item did not get removed")
	}
}

func Test_FIFO_detectLineJumpers(t *testing.T) {
	f := New[testFifoObject](testFifoObjectKeyFunc)

	f.Add(mkFifoObj("foo", 10)) // nolint: errcheck
	f.Add(mkFifoObj("bar", 1))  // nolint: errcheck
	f.Add(mkFifoObj("foo", 11)) // nolint: errcheck
	f.Add(mkFifoObj("foo", 13)) // nolint: errcheck
	f.Add(mkFifoObj("zab", 30)) // nolint: errcheck

	if e, a := 13, Pop[testFifoObject](f).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
	// ensure foo doesn't jump back in line
	f.Add(mkFifoObj("foo", 14)) // nolint: errcheck

	if e, a := 1, Pop[testFifoObject](f).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}

	if e, a := 30, Pop[testFifoObject](f).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}

	if e, a := 14, Pop[testFifoObject](f).val; a != e {
		t.Fatalf("expected %d, got %d", e, a)
	}
}

func Test_FIFO_addIfNotPresent(t *testing.T) {
	f := New[testFifoObject](testFifoObjectKeyFunc)

	f.Add(mkFifoObj("a", 1))             // nolint: errcheck
	f.Add(mkFifoObj("b", 2))             // nolint: errcheck
	f.AddIfNotPresent(mkFifoObj("b", 3)) // nolint: errcheck
	f.AddIfNotPresent(mkFifoObj("c", 4)) // nolint: errcheck

	if e, a := 3, len(f.items); a != e {
		t.Fatalf("expected queue length %d, got %d", e, a)
	}

	expectedValues := []int{1, 2, 4}
	for _, expected := range expectedValues {
		if actual := Pop[testFifoObject](f).val; actual != expected {
			t.Fatalf("expected value %d, got %d", expected, actual)
		}
	}
}

func Test_FIFO_Delete(t *testing.T) {
	f := New(testFifoObjectKeyFunc)

	f.Add(mkFifoObj("a", 1))             // nolint: errcheck
	f.Add(mkFifoObj("b", 2))             // nolint: errcheck
	f.AddIfNotPresent(mkFifoObj("b", 3)) // nolint: errcheck
	f.AddIfNotPresent(mkFifoObj("c", 4)) // nolint: errcheck

	f.Delete(mkFifoObj("b", 2)) // nolint: errcheck
	if e, a := 2, len(f.items); a != e {
		t.Fatalf("expected queue length %d, got %d", e, a)
	}

	expectedValues := []int{1, 4}
	for _, expected := range expectedValues {
		if actual := Pop[testFifoObject](f).val; actual != expected {
			t.Fatalf("expected value %d, got %d", expected, actual)
		}
	}
}

func Test_FIFO_HasSynced(t *testing.T) {
	tests := []struct {
		actions        []func(f *FIFO[testFifoObject])
		expectedSynced bool
	}{
		{
			actions:        []func(f *FIFO[testFifoObject]){},
			expectedSynced: false,
		},
		{
			actions: []func(f *FIFO[testFifoObject]){
				func(f *FIFO[testFifoObject]) {
					f.Add(mkFifoObj("a", 1)) // nolint: errcheck
				},
			},
			expectedSynced: true,
		},
		{
			actions: []func(f *FIFO[testFifoObject]){
				func(f *FIFO[testFifoObject]) {
					f.Replace([]testFifoObject{}, "0") // nolint: errcheck
				},
			},
			expectedSynced: true,
		},
		{
			actions: []func(f *FIFO[testFifoObject]){
				func(f *FIFO[testFifoObject]) {
					f.Replace([]testFifoObject{mkFifoObj("a", 1), mkFifoObj("b", 2)}, "0") // nolint: errcheck
				},
			},
			expectedSynced: false,
		},
		{
			actions: []func(f *FIFO[testFifoObject]){
				func(f *FIFO[testFifoObject]) {
					f.Replace([]testFifoObject{mkFifoObj("a", 1), mkFifoObj("b", 2)}, "0") // nolint: errcheck
				},
				func(f *FIFO[testFifoObject]) { Pop[testFifoObject](f) },
			},
			expectedSynced: false,
		},
		{
			actions: []func(f *FIFO[testFifoObject]){
				func(f *FIFO[testFifoObject]) {
					f.Replace([]testFifoObject{mkFifoObj("a", 1), mkFifoObj("b", 2)}, "0") // nolint: errcheck
				},
				func(f *FIFO[testFifoObject]) { Pop[testFifoObject](f) },
				func(f *FIFO[testFifoObject]) { Pop[testFifoObject](f) },
			},
			expectedSynced: true,
		},
	}

	for i, test := range tests {
		f := New(testFifoObjectKeyFunc)

		for _, action := range test.actions {
			action(f)
		}
		if e, a := test.expectedSynced, f.HasSynced(); a != e {
			t.Errorf("test case %v failed, expected: %v , got %v", i, e, a)
		}
	}
}

// TestFIFO_PopShouldUnblockWhenClosed checks that any blocking Pop on an empty queue
// should unblock and return after Close is called.
func Test_FIFO_PopShouldUnblockWhenClosed(t *testing.T) {
	f := New(testFifoObjectKeyFunc)

	c := make(chan struct{})
	const jobs = 10
	for i := 0; i < jobs; i++ {
		go func() {
			f.Pop(func(obj testFifoObject) error { return nil }) // nolint: errcheck
			c <- struct{}{}
		}()
	}

	runtime.Gosched()
	f.Close()

	if v := f.IsClosed(); !v {
		t.Errorf("test IsClosed failed, expected: %v , got %v", true, v)
	}

	for i := 0; i < jobs; i++ {
		select {
		case <-c:
		case <-time.After(500 * time.Millisecond):
			t.Fatalf("timed out waiting for Pop to return after Close")
		}
	}
}

func Test_FIFO_Resync(t *testing.T) {
	f := New(testFifoObjectKeyFunc)

	if err := f.Resync(); err != nil {
		t.Fatalf("expected value nil, got %d", err)
	}
}
