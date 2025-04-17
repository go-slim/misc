package misc

import (
	"unsafe"
)

// SliceByteToString convert a byte slice to string
// Deprecated: use UnsafeBytesToString func
func SliceByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// UnsafeBytesToString uses Go's unsafe package to convert a byte slice to a string.
func UnsafeBytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// StringToSliceByte convert a string to byte slice.
// Deprecated: use UnsafeStringToBytes func
func StringToSliceByte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// UnsafeStringToBytes uses Go's unsafe package to convert a string to a byte slice.
func UnsafeStringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
