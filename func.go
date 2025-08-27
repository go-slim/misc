package misc

import (
	"reflect"
	"runtime"

	"go-slim.dev/l4g"
)

// Wrap 将一组无参函数组合为一个函数，调用时按顺序执行，遇到错误即返回。
func Wrap(fns ...func() error) func() error {
	return func() error {
		return Call(fns...)
	}
}

// WrapWith 将一组接收相同参数的函数组合为一个函数，调用时按顺序执行，遇到错误即返回。
func WrapWith[T any](fns ...func(T) error) func(val T) error {
	return func(val T) error {
		return CallWith(val, fns...)
	}
}

// Call 依次调用提供的函数，若任一函数返回错误，则立即返回该错误。
func Call(fns ...func() error) error {
	for _, fn := range fns {
		err := fn()
		if err != nil {
			return err
		}
	}
	return nil
}

// CallWith 依次调用提供的函数并传入相同参数，若任一函数返回错误，则立即返回该错误。
func CallWith[T any](val T, fns ...func(T) error) error {
	for _, fn := range fns {
		err := fn(val)
		if err != nil {
			return err
		}
	}
	return nil
}

// MustCall 依次调用函数；若出现错误，打印函数名并以 fatal 退出程序。
func MustCall(fns ...func() error) {
	for _, fn := range fns {
		if err := fn(); err != nil {
			ptr := reflect.ValueOf(fn).Pointer()
			fi := runtime.FuncForPC(ptr)
			l4g.Fatalf("%s failed: %v", fi.Name(), err)
		}
	}
}

// MustCallWith 依次调用函数并传入相同参数；若出现错误，打印函数名并以 fatal 退出程序。
func MustCallWith[T any](val T, fns ...func(T) error) {
	for _, fn := range fns {
		if err := fn(val); err != nil {
			ptr := reflect.ValueOf(fn).Pointer()
			fi := runtime.FuncForPC(ptr)
			l4g.Fatalf("%s(ctx) failed: %v", fi.Name(), err)
		}
	}
}
