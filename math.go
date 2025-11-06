package misc

import (
	"cmp"
)

// MinMax returns the minimum and maximum of two ordered values.
// It returns (min, max) where min is the smaller value and max is the larger value.
// If the values are equal, both return values will be the same.
//
// The function uses generics and works with any type that implements the cmp.Ordered interface,
// including built-in types like int, float64, string, etc.
//
// Example:
//
//	min, max := MinMax(5, 3)        // min=3, max=5
//	min, max := MinMax(1.5, 2.7)    // min=1.5, max=2.7
//	min, max := MinMax("b", "a")    // min="a", max="b"
//	min, max := MinMax(10, 10)      // min=10, max=10 (equal values)
func MinMax[T cmp.Ordered](a, b T) (T, T) {
	switch cmp.Compare(a, b) {
	case -1:
		return a, b // a < b
	case 1:
		return b, a // a > b
	default:
		return a, b // a == b
	}
}

// Clamp limits a value to a specified range [minT, maxT].
// Returns val limited to the range [minT, maxT].
// If val is less than minT, returns minT.
// If val is greater than maxT, returns maxT.
// Otherwise, returns val unchanged.
//
// The function automatically handles cases where minT > maxT by swapping the bounds.
// This makes it more robust and user-friendly.
//
// Example:
//
//	result := Clamp(15, 10, 20)   // result=15 (within range)
//	result := Clamp(5, 10, 20)    // result=10 (below minimum)
//	result := Clamp(25, 10, 20)   // result=20 (above maximum)
//	result := Clamp(15, 20, 10)   // result=15 (bounds automatically swapped)
//	result := Clamp(3.14, 0, 1)   // result=1.0 (float64 example)
func Clamp[T cmp.Ordered](val, minT, maxT T) T {
	a, b := MinMax(minT, maxT)
	return max(min(val, b), a)
}
