package misc

import "reflect"

// IsZero 泛型零值判断
// https://stackoverflow.com/questions/74000242/in-golang-how-to-compare-interface-as-generics-type-to-nil
func IsZero[T any](v T) bool {
	return isZero(reflect.ValueOf(v))
}

func isZero(ref reflect.Value) bool {
	if !ref.IsValid() {
		return true
	}
	if ref.Type().Kind() == reflect.Ptr {
		return isZero(ref.Elem())
	}
	return ref.IsZero()
}
