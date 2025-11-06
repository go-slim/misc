package misc

// Wrap combines a group of parameterless functions into one function, executed in sequence, returning immediately on error.
func Wrap(fns ...func() error) func() error {
	return func() error {
		return Call(fns...)
	}
}

// WrapG combines a group of functions receiving the same parameters into one function, executed in sequence, returning immediately on error.
func WrapG[T any](fns ...func(T) error) func(val T) error {
	return func(val T) error {
		return CallG(val, fns...)
	}
}

// Call calls the provided functions in sequence; if any function returns an error, immediately return that error.
func Call(fns ...func() error) error {
	for _, fn := range fns {
		err := fn()
		if err != nil {
			return err
		}
	}
	return nil
}

// CallG calls the provided functions with the same parameters in sequence; if any function returns an error, immediately return that error.
func CallG[T any](val T, fns ...func(T) error) error {
	for _, fn := range fns {
		err := fn(val)
		if err != nil {
			return err
		}
	}
	return nil
}
