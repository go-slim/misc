package misc

import "testing"

func TestUnsafeBytesStringConversions(t *testing.T) {
	b := []byte("hello")
	s := UnsafeBytesToString(b)
	if s != "hello" {
		t.Fatalf("UnsafeBytesToString mismatch: %q", s)
	}

	bs := UnsafeStringToBytes("world")
	if string(bs) != "world" {
		t.Fatalf("UnsafeStringToBytes mismatch: %q", string(bs))
	}
}

func TestDeprecatedHelpers(t *testing.T) {
	b := []byte("abc")
	s := SliceByteToString(b)
	if s != "abc" {
		t.Fatalf("SliceByteToString mismatch: %q", s)
	}
	bs := StringToSliceByte("xyz")
	if string(bs) != "xyz" {
		t.Fatalf("StringToSliceByte mismatch: %q", string(bs))
	}
}
