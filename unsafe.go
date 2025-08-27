package misc

import (
	"unsafe"
)

// SliceByteToString convert a byte slice to string
// Deprecated: use UnsafeBytesToString func
// SliceByteToString 将字节切片转换为字符串（零拷贝）。
// 注意：属于不安全转换，底层数据共享；若切片内容随后被修改，将影响字符串视图。
func SliceByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// UnsafeBytesToString uses Go's unsafe package to convert a byte slice to a string.
// UnsafeBytesToString 使用 Go 1.20+ 的 unsafe 辅助 API 进行零拷贝转换。
// 注意：结果与原切片共享底层数据，需保证原切片后续不被修改。
func UnsafeBytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// StringToSliceByte convert a string to byte slice.
// Deprecated: use UnsafeStringToBytes func
// StringToSliceByte 将字符串转换为字节切片（零拷贝）。
// 注意：属于不安全转换，返回切片不可修改（否则会触发未定义行为）。
func StringToSliceByte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// UnsafeStringToBytes uses Go's unsafe package to convert a string to a byte slice.
// UnsafeStringToBytes 使用 Go 1.20+ 的 unsafe 辅助 API 将字符串零拷贝转换为字节切片。
// 注意：得到的切片只读，切勿修改其内容。
func UnsafeStringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
