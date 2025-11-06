package misc

import (
	"unsafe"
)

// BytesToString uses Go 1.20+ unsafe helper API to convert byte slice to string with zero copy.
// Note: The result shares underlying data with the original slice, ensure the original slice is not modified subsequently.
func BytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// StringToBytes uses Go 1.20+ unsafe helper API to convert string to byte slice with zero copy.
// Note: The resulting slice is read-only, do not modify its content.
func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
