package misc

import (
	"errors"
	"testing"
)

func TestCall(t *testing.T) {
	var called []int
	err := Call(
		func() error { called = append(called, 1); return nil },
		func() error { called = append(called, 2); return nil },
	)
	if err != nil {
		t.Fatalf("Call returned error: %v", err)
	}
	if len(called) != 2 || called[0] != 1 || called[1] != 2 {
		t.Fatalf("unexpected call order: %v", called)
	}
}

func TestCallStopsOnError(t *testing.T) {
	var called []int
	err := Call(
		func() error { called = append(called, 1); return errors.New("boom") },
		func() error { called = append(called, 2); return nil },
	)
	if err == nil {
		t.Fatal("Call should return the first error")
	}
	if len(called) != 1 || called[0] != 1 {
		t.Fatalf("should stop on first error, called=%v", called)
	}
}

func TestCallWith(t *testing.T) {
	var acc int
	err := CallWith(3,
		func(v int) error { acc += v; return nil },
		func(v int) error { acc += v * 2; return nil },
	)
	if err != nil {
		t.Fatalf("CallWith returned error: %v", err)
	}
	if acc != 3+3*2 {
		t.Fatalf("unexpected acc: %d", acc)
	}
}

func TestWrap(t *testing.T) {
	var x int
	fn := Wrap(
		func() error { x++; return nil },
		func() error { x += 2; return nil },
	)
	if err := fn(); err != nil {
		t.Fatalf("Wrap returned error: %v", err)
	}
	if x != 3 {
		t.Fatalf("unexpected x: %d", x)
	}
}

func TestWrapPropagatesError(t *testing.T) {
	errBoom := errors.New("boom")
	fn := Wrap(
		func() error { return errBoom },
		func() error { t.Fatal("should not be called"); return nil },
	)
	if err := fn(); err != errBoom {
		t.Fatalf("Wrap should propagate first error, got: %v", err)
	}
}

func TestWrapWith(t *testing.T) {
	var s string
	fn := WrapWith(func(v string) error { s = v; return nil }, func(v string) error { s += "!"; return nil })
	if err := fn("hi"); err != nil {
		t.Fatalf("WrapWith returned error: %v", err)
	}
	if s != "hi!" {
		t.Fatalf("unexpected s: %s", s)
	}
}

func TestMustCallSuccess(t *testing.T) {
	// MustCall should not abort when all funcs succeed
	MustCall(func() error { return nil }, func() error { return nil })
}

func TestMustCallWithSuccess(t *testing.T) {
	// MustCallWith should not abort when all funcs succeed
	MustCallWith(1, func(int) error { return nil }, func(int) error { return nil })
}
