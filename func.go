package misc

import (
	"reflect"
	"runtime"

	"zestack.dev/l4g"
)

func Wrap(fns ...func() error) func() error {
	return func() error {
		return Call(fns...)
	}
}

func WrapWith[T any](fns ...func(T) error) func(val T) error {
	return func(val T) error {
		return CallWith(val, fns...)
	}
}

func Call(fns ...func() error) error {
	for _, fn := range fns {
		err := fn()
		if err != nil {
			return err
		}
	}
	return nil
}

func CallWith[T any](val T, fns ...func(T) error) error {
	for _, fn := range fns {
		err := fn(val)
		if err != nil {
			return err
		}
	}
	return nil
}

func MustCall(fns ...func() error) {
	for _, fn := range fns {
		if err := fn(); err != nil {
			ptr := reflect.ValueOf(fn).Pointer()
			fi := runtime.FuncForPC(ptr)
			l4g.Fatalf("%s failed: %v", fi.Name(), err)
		}
	}
}

func MustCallWith[T any](val T, fns ...func(T) error) {
	for _, fn := range fns {
		if err := fn(val); err != nil {
			ptr := reflect.ValueOf(fn).Pointer()
			fi := runtime.FuncForPC(ptr)
			l4g.Fatalf("%s(ctx) failed: %v", fi.Name(), err)
		}
	}
}
