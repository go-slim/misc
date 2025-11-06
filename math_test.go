package misc

import (
	"testing"
)

// Test types for generic testing
type TestOrdered struct {
	value int
}

// Implement comparison for TestOrdered
func (t TestOrdered) Compare(other TestOrdered) int {
	if t.value < other.value {
		return -1
	} else if t.value > other.value {
		return 1
	}
	return 0
}

// Test cases for MinMax function
func TestMinMax(t *testing.T) {
	tests := []struct {
		name    string
		a, b    any
		wantMin any
		wantMax any
	}{
		// Integer tests
		{
			name:    "int a < b",
			a:       3,
			b:       7,
			wantMin: 3,
			wantMax: 7,
		},
		{
			name:    "int a > b",
			a:       10,
			b:       5,
			wantMin: 5,
			wantMax: 10,
		},
		{
			name:    "int a == b",
			a:       42,
			b:       42,
			wantMin: 42,
			wantMax: 42,
		},
		{
			name:    "int negative values",
			a:       -5,
			b:       -2,
			wantMin: -5,
			wantMax: -2,
		},
		{
			name:    "int mixed sign",
			a:       -3,
			b:       5,
			wantMin: -3,
			wantMax: 5,
		},
		// Float tests
		{
			name:    "float64 a < b",
			a:       3.14,
			b:       6.28,
			wantMin: 3.14,
			wantMax: 6.28,
		},
		{
			name:    "float64 a > b",
			a:       9.99,
			b:       1.23,
			wantMin: 1.23,
			wantMax: 9.99,
		},
		{
			name:    "float64 a == b",
			a:       2.718,
			b:       2.718,
			wantMin: 2.718,
			wantMax: 2.718,
		},
		{
			name:    "float32 precision",
			a:       float32(1.1),
			b:       float32(2.2),
			wantMin: float32(1.1),
			wantMax: float32(2.2),
		},
		// String tests
		{
			name:    "string a < b",
			a:       "apple",
			b:       "banana",
			wantMin: "apple",
			wantMax: "banana",
		},
		{
			name:    "string a > b",
			a:       "zebra",
			b:       "yak",
			wantMin: "yak",
			wantMax: "zebra",
		},
		{
			name:    "string a == b",
			a:       "hello",
			b:       "hello",
			wantMin: "hello",
			wantMax: "hello",
		},
		// Uint tests
		{
			name:    "uint a < b",
			a:       uint(5),
			b:       uint(10),
			wantMin: uint(5),
			wantMax: uint(10),
		},
		{
			name:    "uint8 a < b",
			a:       uint8(100),
			b:       uint8(200),
			wantMin: uint8(100),
			wantMax: uint8(200),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use type assertion to handle different types
			switch a := tt.a.(type) {
			case int:
				min, max := MinMax(a, tt.b.(int))
				if min != tt.wantMin.(int) || max != tt.wantMax.(int) {
					t.Errorf("MinMax(%v, %v) = (%v, %v), want (%v, %v)", a, tt.b, min, max, tt.wantMin, tt.wantMax)
				}
			case float64:
				min, max := MinMax(a, tt.b.(float64))
				if min != tt.wantMin.(float64) || max != tt.wantMax.(float64) {
					t.Errorf("MinMax(%v, %v) = (%v, %v), want (%v, %v)", a, tt.b, min, max, tt.wantMin, tt.wantMax)
				}
			case float32:
				min, max := MinMax(a, tt.b.(float32))
				if min != tt.wantMin.(float32) || max != tt.wantMax.(float32) {
					t.Errorf("MinMax(%v, %v) = (%v, %v), want (%v, %v)", a, tt.b, min, max, tt.wantMin, tt.wantMax)
				}
			case string:
				min, max := MinMax(a, tt.b.(string))
				if min != tt.wantMin.(string) || max != tt.wantMax.(string) {
					t.Errorf("MinMax(%v, %v) = (%v, %v), want (%v, %v)", a, tt.b, min, max, tt.wantMin, tt.wantMax)
				}
			case uint:
				min, max := MinMax(a, tt.b.(uint))
				if min != tt.wantMin.(uint) || max != tt.wantMax.(uint) {
					t.Errorf("MinMax(%v, %v) = (%v, %v), want (%v, %v)", a, tt.b, min, max, tt.wantMin, tt.wantMax)
				}
			case uint8:
				min, max := MinMax(a, tt.b.(uint8))
				if min != tt.wantMin.(uint8) || max != tt.wantMax.(uint8) {
					t.Errorf("MinMax(%v, %v) = (%v, %v), want (%v, %v)", a, tt.b, min, max, tt.wantMin, tt.wantMax)
				}
			default:
				t.Errorf("Unsupported type: %T", a)
			}
		})
	}
}

// Test cases for Clamp function
func TestClamp(t *testing.T) {
	tests := []struct {
		name     string
		val      any
		min, max any
		want     any
	}{
		// Integer tests
		{
			name: "int within range",
			val:  15,
			min:  10,
			max:  20,
			want: 15,
		},
		{
			name: "int below minimum",
			val:  5,
			min:  10,
			max:  20,
			want: 10,
		},
		{
			name: "int above maximum",
			val:  25,
			min:  10,
			max:  20,
			want: 20,
		},
		{
			name: "int equal to minimum",
			val:  10,
			min:  10,
			max:  20,
			want: 10,
		},
		{
			name: "int equal to maximum",
			val:  20,
			min:  10,
			max:  20,
			want: 20,
		},
		{
			name: "int min > max (swapped)",
			val:  15,
			min:  20,
			max:  10,
			want: 15,
		},
		{
			name: "int negative values",
			val:  -5,
			min:  -10,
			max:  0,
			want: -5,
		},
		// Float tests
		{
			name: "float64 within range",
			val:  3.14,
			min:  1.0,
			max:  5.0,
			want: 3.14,
		},
		{
			name: "float64 below minimum",
			val:  0.5,
			min:  1.0,
			max:  5.0,
			want: 1.0,
		},
		{
			name: "float64 above maximum",
			val:  6.28,
			min:  1.0,
			max:  5.0,
			want: 5.0,
		},
		{
			name: "float32 precision",
			val:  float32(2.5),
			min:  float32(1.0),
			max:  float32(4.0),
			want: float32(2.5),
		},
		// String tests
		{
			name: "string within range",
			val:  "mango",
			min:  "apple",
			max:  "zebra",
			want: "mango",
		},
		{
			name: "string below minimum",
			val:  "aardvark",
			min:  "apple",
			max:  "zebra",
			want: "apple",
		},
		{
			name: "string above maximum",
			val:  "zzz",
			min:  "apple",
			max:  "zebra",
			want: "zebra",
		},
		{
			name: "string min > max (swapped)",
			val:  "mango",
			min:  "zebra",
			max:  "apple",
			want: "mango",
		},
		// Uint tests
		{
			name: "uint within range",
			val:  uint(50),
			min:  uint(25),
			max:  uint(75),
			want: uint(50),
		},
		{
			name: "uint below minimum",
			val:  uint(10),
			min:  uint(25),
			max:  uint(75),
			want: uint(25),
		},
		{
			name: "uint above maximum",
			val:  uint(100),
			min:  uint(25),
			max:  uint(75),
			want: uint(75),
		},
		// Edge cases
		{
			name: "zero values",
			val:  0,
			min:  0,
			max:  0,
			want: 0,
		},
		{
			name: "same min and max",
			val:  42,
			min:  42,
			max:  42,
			want: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use type assertion to handle different types
			switch val := tt.val.(type) {
			case int:
				min, max := tt.min.(int), tt.max.(int)
				result := Clamp(val, min, max)
				if result != tt.want.(int) {
					t.Errorf("Clamp(%v, %v, %v) = %v, want %v", val, min, max, result, tt.want)
				}
			case float64:
				min, max := tt.min.(float64), tt.max.(float64)
				result := Clamp(val, min, max)
				if result != tt.want.(float64) {
					t.Errorf("Clamp(%v, %v, %v) = %v, want %v", val, min, max, result, tt.want)
				}
			case float32:
				min, max := tt.min.(float32), tt.max.(float32)
				result := Clamp(val, min, max)
				if result != tt.want.(float32) {
					t.Errorf("Clamp(%v, %v, %v) = %v, want %v", val, min, max, result, tt.want)
				}
			case string:
				min, max := tt.min.(string), tt.max.(string)
				result := Clamp(val, min, max)
				if result != tt.want.(string) {
					t.Errorf("Clamp(%v, %v, %v) = %v, want %v", val, min, max, result, tt.want)
				}
			case uint:
				min, max := tt.min.(uint), tt.max.(uint)
				result := Clamp(val, min, max)
				if result != tt.want.(uint) {
					t.Errorf("Clamp(%v, %v, %v) = %v, want %v", val, min, max, result, tt.want)
				}
			default:
				t.Errorf("Unsupported type: %T", val)
			}
		})
	}
}

// Benchmark tests
func BenchmarkMinMax(b *testing.B) {
	b.Run("int", func(b *testing.B) {
		for i := 0; b.Loop(); i++ {
			MinMax(i, i+1)
		}
	})

	b.Run("float64", func(b *testing.B) {
		for i := 0; b.Loop(); i++ {
			MinMax(float64(i), float64(i+1))
		}
	})

	b.Run("string", func(b *testing.B) {
		str1, str2 := "hello", "world"
		for b.Loop() {
			MinMax(str1, str2)
		}
	})
}

func BenchmarkClamp(b *testing.B) {
	b.Run("int", func(b *testing.B) {
		for i := 0; b.Loop(); i++ {
			Clamp(i%100, 10, 90)
		}
	})

	b.Run("float64", func(b *testing.B) {
		for i := 0; b.Loop(); i++ {
			Clamp(float64(i%100), 10.0, 90.0)
		}
	})

	b.Run("string", func(b *testing.B) {
		vals := []string{"apple", "banana", "cherry", "date", "elderberry"}
		for i := 0; b.Loop(); i++ {
			Clamp(vals[i%len(vals)], "banana", "date")
		}
	})
}

// Example tests for documentation
func ExampleMinMax() {
	// Integer example
	min, max := MinMax(5, 3)
	println(min, max)
	// Output: 3 5
}

func ExampleMinMax_int() {
	// Integer example
	min, max := MinMax(5, 3)
	println(min, max)
	// Output: 3 5
}

func ExampleMinMax_float64() {
	// Float example
	min, max := MinMax(1.5, 2.7)
	println(min, max)
	// Output: 1.5 2.7
}

func ExampleMinMax_string() {
	// String example
	min, max := MinMax("b", "a")
	println(min, max)
	// Output: a b
}

func ExampleClamp() {
	// Within range
	result := Clamp(15, 10, 20)
	println(result)
	// Below minimum
	result = Clamp(5, 10, 20)
	println(result)
	// Above maximum
	result = Clamp(25, 10, 20)
	println(result)
	// Bounds swapped (automatically handled)
	result = Clamp(15, 20, 10)
	println(result)
	// Output: 15
	// 10
	// 20
	// 15
}
