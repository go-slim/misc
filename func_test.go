package misc

import (
	"errors"
	"fmt"
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

func TestCallG(t *testing.T) {
	var acc int
	err := CallG(3,
		func(v int) error { acc += v; return nil },
		func(v int) error { acc += v * 2; return nil },
	)
	if err != nil {
		t.Fatalf("CallG returned error: %v", err)
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

func TestWrapG(t *testing.T) {
	var s string
	fn := WrapG(func(v string) error { s = v; return nil }, func(v string) error { s += "!"; return nil })
	if err := fn("hi"); err != nil {
		t.Fatalf("WrapG returned error: %v", err)
	}
	if s != "hi!" {
		t.Fatalf("unexpected s: %s", s)
	}
}

// Example tests for godoc

func ExampleCall() {
	// Execute multiple functions in sequence
	err := Call(
		func() error {
			fmt.Println("Step 1")
			return nil
		},
		func() error {
			fmt.Println("Step 2")
			return nil
		},
		func() error {
			fmt.Println("Step 3")
			return nil
		},
	)
	if err != nil {
		fmt.Println("Error:", err)
	}
	// Output: Step 1
	// Step 2
	// Step 3
}

func ExampleCall_errorHandling() {
	// Stops on first error
	err := Call(
		func() error {
			fmt.Println("Step 1")
			return nil
		},
		func() error {
			return errors.New("something went wrong")
		},
		func() error {
			fmt.Println("Step 3 - not executed")
			return nil
		},
	)
	fmt.Println("Error:", err)
	// Output: Step 1
	// Error: something went wrong
}

func ExampleCallG() {
	// Execute multiple functions with the same parameter
	value := 10
	err := CallG(value,
		func(v int) error {
			fmt.Println("Value:", v)
			return nil
		},
		func(v int) error {
			fmt.Println("Doubled:", v*2)
			return nil
		},
		func(v int) error {
			fmt.Println("Tripled:", v*3)
			return nil
		},
	)
	if err != nil {
		fmt.Println("Error:", err)
	}
	// Output: Value: 10
	// Doubled: 20
	// Tripled: 30
}

func ExampleWrap() {
	// Create a composite function
	var results []int
	composedFunc := Wrap(
		func() error {
			results = append(results, 1)
			return nil
		},
		func() error {
			results = append(results, 2)
			return nil
		},
		func() error {
			results = append(results, 3)
			return nil
		},
	)

	// Execute the composed function
	_ = composedFunc()
	fmt.Println("First call:", results)

	// Execute again
	_ = composedFunc()
	fmt.Println("Second call:", results)
	// Output: First call: [1 2 3]
	// Second call: [1 2 3 1 2 3]
}

func ExampleWrapG() {
	// Create a composite function with parameters
	processString := WrapG(
		func(s string) error {
			fmt.Println("Processing:", s)
			return nil
		},
		func(s string) error {
			fmt.Println("Length:", len(s))
			return nil
		},
	)

	_ = processString("Hello")
	_ = processString("World!")
	// Output: Processing: Hello
	// Length: 5
	// Processing: World!
	// Length: 6
}
