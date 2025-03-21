package misc

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
