package misc

import (
	"fmt"
	"reflect"
	"testing"
)

// Test cases for Zero function
func TestZero(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func() any
		expected any
	}{
		{
			name:     "int zero",
			testFunc: func() any { return Zero[int]() },
			expected: 0,
		},
		{
			name:     "int8 zero",
			testFunc: func() any { return Zero[int8]() },
			expected: int8(0),
		},
		{
			name:     "int16 zero",
			testFunc: func() any { return Zero[int16]() },
			expected: int16(0),
		},
		{
			name:     "int32 zero",
			testFunc: func() any { return Zero[int32]() },
			expected: int32(0),
		},
		{
			name:     "int64 zero",
			testFunc: func() any { return Zero[int64]() },
			expected: int64(0),
		},
		{
			name:     "uint zero",
			testFunc: func() any { return Zero[uint]() },
			expected: uint(0),
		},
		{
			name:     "uint8 zero",
			testFunc: func() any { return Zero[uint8]() },
			expected: uint8(0),
		},
		{
			name:     "float32 zero",
			testFunc: func() any { return Zero[float32]() },
			expected: float32(0),
		},
		{
			name:     "float64 zero",
			testFunc: func() any { return Zero[float64]() },
			expected: float64(0),
		},
		{
			name:     "string zero",
			testFunc: func() any { return Zero[string]() },
			expected: "",
		},
		{
			name:     "bool zero",
			testFunc: func() any { return Zero[bool]() },
			expected: false,
		},
		{
			name:     "pointer zero",
			testFunc: func() any { return Zero[*int]() },
			expected: (*int)(nil),
		},
		{
			name:     "slice zero",
			testFunc: func() any { return Zero[[]int]() },
			expected: []int(nil),
		},
		{
			name:     "map zero",
			testFunc: func() any { return Zero[map[string]int]() },
			expected: map[string]int(nil),
		},
		{
			name:     "channel zero",
			testFunc: func() any { return Zero[chan int]() },
			expected: chan int(nil),
		},
		{
			name:     "interface zero",
			testFunc: func() any { return Zero[any]() },
			expected: nil,
		},
		{
			name:     "function zero",
			testFunc: func() any { return Zero[func()]() },
			expected: (func())(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.testFunc()
			// Handle incomparable types like slices, maps, functions
			switch tt.name {
			case "slice zero", "map zero", "channel zero", "function zero":
				if result == nil && tt.expected == nil {
					return // success
				}
				if (result == nil) != (tt.expected == nil) {
					t.Errorf("Zero() = %v, want %v", result, tt.expected)
				}
			default:
				// For comparable types, use direct comparison
				if result != tt.expected {
					t.Errorf("Zero() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

// Test cases for Ptr function
func TestPtr(t *testing.T) {
	tests := []struct {
		name string
		x    interface{}
	}{
		{
			name: "int pointer",
			x:    42,
		},
		{
			name: "string pointer",
			x:    "hello world",
		},
		{
			name: "bool pointer",
			x:    true,
		},
		{
			name: "float pointer",
			x:    3.14159,
		},
		{
			name: "struct pointer",
			x:    struct{ Name string }{Name: "test"},
		},
		{
			name: "slice pointer",
			x:    []int{1, 2, 3},
		},
		{
			name: "map pointer",
			x:    map[string]int{"key": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use type assertion to test different types
			switch v := tt.x.(type) {
			case int:
				ptr := Ptr(v)
				if ptr == nil {
					t.Fatal("Ptr() returned nil")
				}
				if *ptr != v {
					t.Errorf("Ptr(%v) = %v, want %v", v, *ptr, v)
				}
			case string:
				ptr := Ptr(v)
				if ptr == nil {
					t.Fatal("Ptr() returned nil")
				}
				if *ptr != v {
					t.Errorf("Ptr(%v) = %v, want %v", v, *ptr, v)
				}
			case bool:
				ptr := Ptr(v)
				if ptr == nil {
					t.Fatal("Ptr() returned nil")
				}
				if *ptr != v {
					t.Errorf("Ptr(%v) = %v, want %v", v, *ptr, v)
				}
			case float64:
				ptr := Ptr(v)
				if ptr == nil {
					t.Fatal("Ptr() returned nil")
				}
				if *ptr != v {
					t.Errorf("Ptr(%v) = %v, want %v", v, *ptr, v)
				}
			case []int:
				ptr := Ptr(v)
				if ptr == nil {
					t.Fatal("Ptr() returned nil")
				}
				// For slices, just check length and elements
				if len(*ptr) != len(v) {
					t.Errorf("Ptr(%v) slice length = %d, want %d", v, len(*ptr), len(v))
				}
				for i := range *ptr {
					if (*ptr)[i] != v[i] {
						t.Errorf("Ptr(%v) slice element [%d] = %v, want %v", v, i, (*ptr)[i], v[i])
					}
				}
			case map[string]int:
				ptr := Ptr(v)
				if ptr == nil {
					t.Fatal("Ptr() returned nil")
				}
				// For maps, check length and elements
				if len(*ptr) != len(v) {
					t.Errorf("Ptr(%v) map length = %d, want %d", v, len(*ptr), len(v))
				}
				for k, val := range *ptr {
					if v[k] != val {
						t.Errorf("Ptr(%v) map key %s = %v, want %v", v, k, val, v[k])
					}
				}
			}
		})
	}
}

// Test cases for Nil function
func TestNil(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func() any
	}{
		{
			name:     "int pointer nil",
			testFunc: func() any { return Nil[int]() },
		},
		{
			name:     "string pointer nil",
			testFunc: func() any { return Nil[string]() },
		},
		{
			name:     "bool pointer nil",
			testFunc: func() any { return Nil[bool]() },
		},
		{
			name:     "struct pointer nil",
			testFunc: func() any { return Nil[struct{ Name string }]() },
		},
		{
			name:     "slice pointer nil",
			testFunc: func() any { return Nil[[]int]() },
		},
		{
			name:     "map pointer nil",
			testFunc: func() any { return Nil[map[string]int]() },
		},
		{
			name:     "function pointer nil",
			testFunc: func() any { return Nil[func()]() },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.testFunc()
			// For typed nil pointers, check if they are nil
			// Use reflection to properly check nil status
			val := reflect.ValueOf(result)
			if val.Kind() == reflect.Pointer {
				if !val.IsNil() {
					t.Errorf("Nil() = %v, want nil", result)
				}
			} else {
				// For non-pointer types, compare directly
				if result != nil {
					t.Errorf("Nil() = %v, want nil", result)
				}
			}
		})
	}
}

// Test cases for IsZero function
func TestIsZero(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		// Zero values
		{
			name:     "int zero",
			value:    0,
			expected: true,
		},
		{
			name:     "int8 zero",
			value:    int8(0),
			expected: true,
		},
		{
			name:     "uint zero",
			value:    uint(0),
			expected: true,
		},
		{
			name:     "float32 zero",
			value:    float32(0),
			expected: true,
		},
		{
			name:     "float64 zero",
			value:    0.0,
			expected: true,
		},
		{
			name:     "string zero",
			value:    "",
			expected: true,
		},
		{
			name:     "bool zero",
			value:    false,
			expected: true,
		},
		{
			name:     "nil pointer",
			value:    (*int)(nil),
			expected: true,
		},
		{
			name:     "nil slice",
			value:    []int(nil),
			expected: true,
		},
		{
			name:     "nil map",
			value:    map[string]int(nil),
			expected: true,
		},
		{
			name:     "empty struct",
			value:    struct{}{},
			expected: true,
		},
		{
			name:     "nil nested pointer",
			value:    (**int)(nil),
			expected: true,
		},
		// Non-zero values
		{
			name:     "int non-zero",
			value:    42,
			expected: false,
		},
		{
			name:     "uint non-zero",
			value:    uint(42),
			expected: false,
		},
		{
			name:     "float non-zero",
			value:    3.14,
			expected: false,
		},
		{
			name:     "string non-zero",
			value:    "hello",
			expected: false,
		},
		{
			name:     "bool non-zero",
			value:    true,
			expected: false,
		},
		{
			name:     "non-nil pointer",
			value:    Ptr(42),
			expected: false,
		},
		{
			name:     "non-nil slice",
			value:    []int{},
			expected: false,
		},
		{
			name:     "non-nil map",
			value:    map[string]int{},
			expected: false,
		},
		{
			name:     "non-empty struct",
			value:    struct{ Value int }{Value: 1},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use type assertion to test different types
			switch v := tt.value.(type) {
			case int:
				result := IsZero(v)
				if result != tt.expected {
					t.Errorf("IsZero(%v) = %v, want %v", v, result, tt.expected)
				}
			case uint:
				result := IsZero(v)
				if result != tt.expected {
					t.Errorf("IsZero(%v) = %v, want %v", v, result, tt.expected)
				}
			case float64:
				result := IsZero(v)
				if result != tt.expected {
					t.Errorf("IsZero(%v) = %v, want %v", v, result, tt.expected)
				}
			case string:
				result := IsZero(v)
				if result != tt.expected {
					t.Errorf("IsZero(%v) = %v, want %v", v, result, tt.expected)
				}
			case bool:
				result := IsZero(v)
				if result != tt.expected {
					t.Errorf("IsZero(%v) = %v, want %v", v, result, tt.expected)
				}
			case int8:
				result := IsZero(v)
				if result != tt.expected {
					t.Errorf("IsZero(%v) = %v, want %v", v, result, tt.expected)
				}
			case float32:
				result := IsZero(v)
				if result != tt.expected {
					t.Errorf("IsZero(%v) = %v, want %v", v, result, tt.expected)
				}
			case struct{}:
				result := IsZero(v)
				if result != tt.expected {
					t.Errorf("IsZero(%v) = %v, want %v", v, result, tt.expected)
				}
			case struct{ Value int }:
				result := IsZero(v)
				if result != tt.expected {
					t.Errorf("IsZero(%v) = %v, want %v", v, result, tt.expected)
				}
			case *int:
				result := IsZero(v)
				if result != tt.expected {
					t.Errorf("IsZero(%v) = %v, want %v", v, result, tt.expected)
				}
			case []int:
				result := IsZero(v)
				if result != tt.expected {
					t.Errorf("IsZero(%v) = %v, want %v", v, result, tt.expected)
				}
			case map[string]int:
				result := IsZero(v)
				if result != tt.expected {
					t.Errorf("IsZero(%v) = %v, want %v", v, result, tt.expected)
				}
			case **int:
				result := IsZero(v)
				if result != tt.expected {
					t.Errorf("IsZero(%v) = %v, want %v", v, result, tt.expected)
				}
			default:
				t.Errorf("Unsupported type: %T", v)
			}
		})
	}
}

// Test cases for IsNil function
func TestIsNil(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "nil interface",
			value:    nil,
			expected: true,
		},
		{
			name:     "nil pointer",
			value:    (*int)(nil),
			expected: true,
		},
		{
			name:     "nil slice",
			value:    []int(nil),
			expected: true,
		},
		{
			name:     "nil map",
			value:    map[string]int(nil),
			expected: true,
		},
		{
			name:     "nil channel",
			value:    (<-chan int)(nil),
			expected: true,
		},
		{
			name:     "nil function",
			value:    (func())(nil),
			expected: true,
		},
		{
			name:     "non-nil interface with nil value",
			value:    any((*int)(nil)),
			expected: true,
		},
		{
			name:     "non-nil pointer",
			value:    Ptr(42),
			expected: false,
		},
		{
			name:     "non-nil slice",
			value:    []int{},
			expected: false,
		},
		{
			name:     "non-nil map",
			value:    map[string]int{},
			expected: false,
		},
		{
			name:     "non-nil channel",
			value:    make(chan int),
			expected: false,
		},
		{
			name:     "int value",
			value:    42,
			expected: false,
		},
		{
			name:     "string value",
			value:    "hello",
			expected: false,
		},
		{
			name:     "bool value",
			value:    true,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNil(tt.value)
			if result != tt.expected {
				t.Errorf("IsNil(%v) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

// Test cases for Coalesce function
func TestCoalesce(t *testing.T) {
	tests := []struct {
		name     string
		values   []interface{}
		testFunc func() any
		expected interface{}
	}{
		{
			name:     "first non-zero string",
			values:   []interface{}{"", "default", "fallback"},
			testFunc: func() any { return Coalesce("", "default", "fallback") },
			expected: "default",
		},
		{
			name:     "first non-zero int",
			values:   []interface{}{0, 42, 100},
			testFunc: func() any { return Coalesce(0, 42, 100) },
			expected: 42,
		},
		{
			name:     "first non-zero bool",
			values:   []interface{}{false, true},
			testFunc: func() any { return Coalesce(false, true) },
			expected: true,
		},
		{
			name:     "all zero values strings",
			values:   []interface{}{"", "", ""},
			testFunc: func() any { return Coalesce("", "", "") },
			expected: "",
		},
		{
			name:     "all zero values ints",
			values:   []interface{}{0, 0, 0},
			testFunc: func() any { return Coalesce(0, 0, 0) },
			expected: 0,
		},
		{
			name:     "single value",
			values:   []interface{}{"single"},
			testFunc: func() any { return Coalesce("single") },
			expected: "single",
		},
		{
			name:     "no values",
			values:   []interface{}{},
			testFunc: func() any { return Coalesce[string]() },
			expected: "",
		},
		{
			name:     "mixed types strings",
			values:   []interface{}{"", "value", ""},
			testFunc: func() any { return Coalesce("", "value", "") },
			expected: "value",
		},
		{
			name:     "mixed types ints",
			values:   []interface{}{0, 0, 5},
			testFunc: func() any { return Coalesce(0, 0, 5) },
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.testFunc()
			if result != tt.expected {
				t.Errorf("Coalesce(%v) = %v, want %v", tt.values, result, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkZero(b *testing.B) {
	b.Run("int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Zero[int]()
		}
	})
	b.Run("string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Zero[string]()
		}
	})
	b.Run("pointer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Zero[*int]()
		}
	})
}

func BenchmarkPtr(b *testing.B) {
	b.Run("int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Ptr(i)
		}
	})
	b.Run("string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Ptr("test")
		}
	})
	b.Run("struct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Ptr(struct{ Val int }{Val: i})
		}
	})
}

func BenchmarkIsNil(b *testing.B) {
	b.Run("nil pointer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = IsNil((*int)(nil))
		}
	})
	b.Run("non-nil pointer", func(b *testing.B) {
		ptr := Ptr(42)
		for i := 0; i < b.N; i++ {
			_ = IsNil(ptr)
		}
	})
	b.Run("mixed", func(b *testing.B) {
		values := []any{nil, Ptr(1), "", "test"}
		for i := 0; i < b.N; i++ {
			_ = IsNil(values[i%len(values)])
		}
	})
}

func BenchmarkCoalesce(b *testing.B) {
	b.Run("int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Coalesce(0, i%2, 42)
		}
	})
	b.Run("string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Coalesce("", "default", "fallback")
		}
	})
}

// Example tests for documentation
func ExampleZero() {
	// Integer zero
	zeroInt := Zero[int]()
	fmt.Println(zeroInt)
	// Output: 0
}

func ExampleZero_string() {
	// String zero
	zeroStr := Zero[string]()
	fmt.Print(zeroStr)
	// Output:
}

func ExampleZero_bool() {
	// Boolean zero
	zeroBool := Zero[bool]()
	fmt.Println(zeroBool)
	// Output: false
}

func ExamplePtr() {
	// Create pointer to integer
	ptr := Ptr(42)
	fmt.Println(*ptr)
	// Output: 42
}

func ExamplePtr_string() {
	// Create pointer to string
	strPtr := Ptr("hello world")
	fmt.Println(*strPtr)
	// Output: hello world
}

func ExampleNil() {
	// Create nil pointer
	var intPtr *int = Nil[int]()
	if intPtr == nil {
		fmt.Println("pointer is nil")
	}
	// Output: pointer is nil
}

func ExampleIsZero() {
	// Check zero values
	fmt.Println(IsZero(0))     // true
	fmt.Println(IsZero(""))    // true
	fmt.Println(IsZero(false)) // true
	fmt.Println(IsZero(42))    // false
	// Output: true
	// true
	// true
	// false
}

func ExampleIsNil() {
	// Check nil values
	fmt.Println(IsNil(nil))            // true
	fmt.Println(IsNil((*int)(nil)))    // true
	fmt.Println(IsNil([]int(nil)))     // true
	fmt.Println(IsNil(make(chan int))) // false
	fmt.Println(IsNil(42))             // false
	// Output: true
	// true
	// true
	// false
	// false
}

func ExampleCoalesce() {
	// Return first non-zero value
	result := Coalesce("", "default", "fallback")
	fmt.Println(result)
	// Output: default
}

func ExampleCoalesce_int() {
	// Return first non-zero int
	result := Coalesce(0, 42, 100)
	fmt.Println(result)
	// Output: 42
}

// Additional IsZero tests moved from validate_test.go
type myStruct struct{ A int }

func TestIsZero_Basics(t *testing.T) {
	if !IsZero(0) || !IsZero(0.0) || !IsZero("") || !IsZero(false) {
		t.Fatal("basic zero checks failed")
	}
	if IsZero(1) || IsZero("x") || IsZero(true) {
		t.Fatal("basic non-zero checks failed")
	}
}

func TestIsZero_PointersAndStructs(t *testing.T) {
	var p *int
	if !IsZero(p) {
		t.Fatal("nil pointer should be zero")
	}
	n := 0
	pp := &p
	if !IsZero(pp) { // pointer to nil pointer -> zero
		t.Fatal("pointer to nil pointer should be zero")
	}
	p = &n
	if !IsZero(p) { // pointer to zero value int -> still zero
		t.Fatal("pointer to zero value should be zero")
	}
	m := myStruct{}
	if !IsZero(m) || IsZero(myStruct{A: 1}) {
		t.Fatal("struct zero checks failed")
	}
}

// Benchmark tests moved from bench_test.go
func BenchmarkIsZero(b *testing.B) {
	type T struct{ N int }
	cases := []any{0, "", (*T)(nil), T{}, []int(nil)}
	b.ReportAllocs()
	for b.Loop() {
		for _, v := range cases {
			_ = IsZero(v)
		}
	}
}
