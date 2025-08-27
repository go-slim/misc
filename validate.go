package misc

import "reflect"

// IsZero 泛型零值判断
// https://stackoverflow.com/questions/74000242/in-golang-how-to-compare-interface-as-generics-type-to-nil
// IsZero 判断任意泛型值是否为零值。
// 若为指针类型，会递归地解引用直至基础值再判断。
func IsZero[T any](v T) bool {
	return isZero(reflect.ValueOf(v))
}

// isZero 对 reflect.Value 执行零值判定。
// 无效值（如 nil）视为零值；指针将被递归解引用。
func isZero(ref reflect.Value) bool {
	if !ref.IsValid() {
		return true
	}
	if ref.Type().Kind() == reflect.Ptr {
		return isZero(ref.Elem())
	}
	return ref.IsZero()
}
