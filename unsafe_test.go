package misc

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

func TestBytesToString(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "empty slice",
			input:    []byte{},
			expected: "",
		},
		{
			name:     "nil slice",
			input:    nil,
			expected: "",
		},
		{
			name:     "simple text",
			input:    []byte("hello world"),
			expected: "hello world",
		},
		{
			name:     "unicode text",
			input:    []byte("你好世界"),
			expected: "你好世界",
		},
		{
			name:     "mixed content",
			input:    []byte("Hello 世界 123!"),
			expected: "Hello 世界 123!",
		},
		{
			name:     "single byte",
			input:    []byte{'A'},
			expected: "A",
		},
		{
			name:     "zero bytes",
			input:    []byte{0, 0, 0},
			expected: string([]byte{0, 0, 0}),
		},
		{
			name:     "long text",
			input:    []byte("this is a longer text to test the bytes to string conversion functionality with more content"),
			expected: "this is a longer text to test the bytes to string conversion functionality with more content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BytesToString(tt.input)
			if result != tt.expected {
				t.Errorf("BytesToString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestStringToBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []byte{},
		},
		{
			name:     "simple text",
			input:    "hello world",
			expected: []byte("hello world"),
		},
		{
			name:     "unicode text",
			input:    "你好世界",
			expected: []byte("你好世界"),
		},
		{
			name:     "mixed content",
			input:    "Hello 世界 123!",
			expected: []byte("Hello 世界 123!"),
		},
		{
			name:     "single character",
			input:    "A",
			expected: []byte{'A'},
		},
		{
			name:     "zero characters",
			input:    string([]byte{0, 0, 0}),
			expected: []byte{0, 0, 0},
		},
		{
			name:     "long text",
			input:    "this is a longer text to test the string to bytes conversion functionality with more content",
			expected: []byte("this is a longer text to test the string to bytes conversion functionality with more content"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringToBytes(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("StringToBytes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestZeroCopyGuarantee(t *testing.T) {
	// Test that BytesToString shares underlying data
	original := []byte("test data")
	str := BytesToString(original)

	// Change original slice
	original[0] = 'T'

	// String should reflect the change (zero copy property)
	if str != "Test data" {
		t.Errorf("Zero copy guarantee failed: string should reflect slice changes")
	}
}

func TestReadOnlyGuarantee(t *testing.T) {
	// Test that StringToBytes creates read-only slice
	str := "test data"
	bytes := StringToBytes(str)

	// Verify that the string header and slice header point to same data
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))
	bytesHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))

	if strHeader.Data != bytesHeader.Data {
		t.Errorf("StringToBytes should point to same underlying data")
	}

	if strHeader.Len != bytesHeader.Len {
		t.Errorf("Length mismatch: string len=%d, slice len=%d", strHeader.Len, bytesHeader.Len)
	}
}

func TestRoundTrip(t *testing.T) {
	testStrings := []string{
		"",
		"hello",
		"你好世界",
		"Hello 世界 123!",
		string([]byte{0, 1, 2, 255}),
	}

	for _, original := range testStrings {
		t.Run(original, func(t *testing.T) {
			// String -> bytes -> string
			bytes := StringToBytes(original)
			result := BytesToString(bytes)

			if result != original {
				t.Errorf("Round trip failed: original=%q, result=%q", original, result)
			}
		})
	}

	testBytes := [][]byte{
		nil,
		{},
		[]byte("hello"),
		[]byte("你好世界"),
		[]byte("Hello 世界 123!"),
		{0, 1, 2, 255},
	}

	for _, original := range testBytes {
		t.Run(string(original), func(t *testing.T) {
			// Bytes -> string -> bytes
			str := BytesToString(original)
			result := StringToBytes(str)

			if !reflect.DeepEqual(result, original) {
				t.Errorf("Round trip failed: original=%v, result=%v", original, result)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("large slice", func(t *testing.T) {
		// Test with larger data
		data := make([]byte, 10000)
		for i := range data {
			data[i] = byte(i % 256)
		}

		str := BytesToString(data)
		result := StringToBytes(str)

		if !reflect.DeepEqual(result, data) {
			t.Errorf("Large data round trip failed")
		}
	})

	t.Run("string with special characters", func(t *testing.T) {
		specialStr := "Special chars: \n\t\r\"'\\"
		bytes := StringToBytes(specialStr)
		result := BytesToString(bytes)

		if result != specialStr {
			t.Errorf("Special characters not preserved: original=%q, result=%q", specialStr, result)
		}
	})
}

// Benchmarks
func BenchmarkBytesToString(b *testing.B) {
	data := []byte("this is a test string for benchmarking")
	b.ResetTimer()
	for b.Loop() {
		_ = BytesToString(data)
	}
}

func BenchmarkStringToBytes(b *testing.B) {
	str := "this is a test string for benchmarking"
	b.ResetTimer()
	for b.Loop() {
		_ = StringToBytes(str)
	}
}

func BenchmarkStringConversionStandard(b *testing.B) {
	// Standard conversion for comparison
	data := []byte("this is a test string for benchmarking")
	b.ResetTimer()
	for b.Loop() {
		_ = string(data)
	}
}

func BenchmarkByteSliceConversionStandard(b *testing.B) {
	// Standard conversion for comparison
	str := "this is a test string for benchmarking"
	b.ResetTimer()
	for b.Loop() {
		_ = []byte(str)
	}
}

// Benchmark tests moved from bench_test.go
func BenchmarkUnsafeConversions(b *testing.B) {
	b.Run("bytes_to_string", func(b *testing.B) {
		bs := []byte(strings.Repeat("a", 64))
		b.ReportAllocs()
		for b.Loop() {
			_ = BytesToString(bs)
		}
	})

	b.Run("string_to_bytes", func(b *testing.B) {
		s := strings.Repeat("a", 64)
		b.ReportAllocs()
		for b.Loop() {
			_ = StringToBytes(s)
		}
	})
}

// Example tests for godoc

func ExampleBytesToString() {
	data := []byte("Hello, World!")
	str := BytesToString(data)
	fmt.Println(str)
	fmt.Println("Type:", reflect.TypeOf(str))

	// WARNING: Do not modify the original slice after conversion
	// as they share the same underlying memory
	// Output: Hello, World!
	// Type: string
}

func ExampleStringToBytes() {
	str := "Hello, Go!"
	bytes := StringToBytes(str)
	fmt.Println(string(bytes))
	fmt.Println("Length:", len(bytes))

	// WARNING: The returned slice is READ-ONLY
	// Modifying it will cause undefined behavior
	// Output: Hello, Go!
	// Length: 10
}

func ExampleBytesToString_performance() {
	// Zero-copy conversion is useful for performance-critical code
	largeData := make([]byte, 1000)
	for i := range largeData {
		largeData[i] = 'A'
	}

	// No memory allocation for the conversion itself
	str := BytesToString(largeData)
	fmt.Println("Converted:", len(str), "bytes")

	// Output: Converted: 1000 bytes
}

func ExampleStringToBytes_readonly() {
	// Useful when you need to pass string data to APIs that accept []byte
	// but you don't want to allocate a new slice
	text := "Read-only data"
	bytes := StringToBytes(text)

	// Safe read operations
	fmt.Println("First byte:", bytes[0])
	fmt.Println("Length:", len(bytes))

	// UNSAFE: Do NOT modify the bytes
	// bytes[0] = 'X' // This would cause undefined behavior!

	// Output: First byte: 82
	// Length: 14
}
