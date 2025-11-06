// Package misc provides utility functions for generic types and validation.
package misc

import "reflect"

// Zero returns the zero value for type T.
// The zero value is the default value for any type in Go:
//   - 0 for numeric types
//   - false for boolean
//   - "" for strings
//   - nil for pointers, slices, maps, channels, interfaces, and functions
//
// This function is useful when you need a zero value for a generic type.
//
// Example:
//
//	zeroInt := Zero[int]()        // zeroInt = 0
//	zeroStr := Zero[string]()    // zeroStr = ""
//	zeroPtr := Zero[*int]()      // zeroPtr = nil
//	zeroBool := Zero[bool]()     // zeroBool = false
func Zero[T any]() T {
	var t T
	return t
}

// Ptr returns a pointer copy of value.
// It takes a value of any type and returns a pointer to a copy of that value.
// This is useful when you need to create a pointer to a value literal
// or convert a value to a pointer type.
//
// Note: For reference types (slices, maps, channels, etc.),
// this creates a pointer to the same underlying reference.
//
// Example:
//
//	ptr := Ptr(42)               // ptr is *int pointing to value 42
//	strPtr := Ptr("hello")       // strPtr is *string pointing to "hello"
//	arrPtr := Ptr([]int{1, 2, 3}) // arrPtr is *[]int pointing to the same slice
func Ptr[T any](x T) *T {
	return &x
}

// Nil returns a nil pointer of type T.
// This is useful when you need to explicitly return or use a nil pointer
// for a specific generic type, especially in generic functions or interfaces.
//
// Example:
//
//	var intPtr *int = Nil[int]()      // intPtr is nil
//	var strPtr *string = Nil[string]() // strPtr is nil
//	var slicePtr *[]int = Nil[[]int]()  // slicePtr is nil
func Nil[T any]() *T {
	return nil
}

// IsZero determines if any generic value is zero.
// It uses reflection to perform deep zero value checking, recursively
// dereferencing pointers until the base value for judgment.
//
// This function handles complex types better than simple DeepEqual comparison,
// especially for nested structs and custom types.
//
// Example:
//
//	IsZero(0)                    // true
//	IsZero("")                   // true
//	IsZero(false)                // true
//	IsZero((*int)(nil))          // true (nil pointer)
//	IsZero(42)                   // false
//	IsZero("hello")              // false
func IsZero[T any](v T) bool {
	return isZero(reflect.ValueOf(v))
}

// isZero performs zero value determination on reflect.Value.
// Invalid values (such as nil) are considered zero values; pointers will be recursively dereferenced.
func isZero(ref reflect.Value) bool {
	if !ref.IsValid() {
		return true
	}
	if ref.Type().Kind() == reflect.Pointer {
		return isZero(ref.Elem())
	}
	return ref.IsZero()
}

// IsNil checks if a value is nil or if it's a reference type with a nil underlying value.
// This function works with any value and returns true for:
//   - nil interface{}
//   - nil pointers
//   - nil slices
//   - nil maps
//   - nil channels
//   - nil functions
//   - nil unsafe pointers
//
// Example:
//
//	IsNil(nil)                    // true
//	IsNil((*int)(nil))            // true
//	IsNil([]int(nil))             // true
//	IsNil(map[string]int(nil))    // true
//	IsNil(make(chan int))          // false
//	IsNil(42)                     // false
//	IsNil("hello")                 // false
func IsNil(x any) bool {
	if x == nil {
		return true
	}
	v := reflect.ValueOf(x)
	switch v.Kind() { //nolint:exhaustive
	case reflect.Chan,
		reflect.Func,
		reflect.Map,
		reflect.Pointer,
		reflect.UnsafePointer,
		reflect.Interface,
		reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

// Coalesce returns the first non-zero value from the provided arguments.
// It iterates through the values in order and returns the first one that
// is not the zero value for its type. If all values are zero values,
// it returns the zero value for type T.
//
// This function is useful for providing fallback values or defaults.
//
// Example:
//
//	result := Coalesce("", "default", "fallback") // result = "default"
//	result := Coalesce(0, 42, 100)               // result = 42
//	result := Coalesce("", "", "value")           // result = "value"
//	result := Coalesce[int]()                    // result = 0 (all zero)
func Coalesce[T any](values ...T) T {
	zero := Zero[T]()
	for _, v := range values {
		if !IsZero(v) {
			return v
		}
	}
	return zero
}
