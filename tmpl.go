package misc

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// TagFunc can be used as a substitution value in the map passed to Interpolate.
// Interpolate functions pass tag (placeholder) name in 'tag' argument.
//
// TagFunc must be safe to call from concurrently running goroutines.
//
// TagFunc must write contents to w and return the number of bytes written.
type TagFunc func(w io.Writer, tag string) (int, error)

// Strtr translate characters or replace substrings
func Strtr(template string, data map[string]any) (string, error) {
	return Tmpl(template, "{", "}", data)
}

func Tmpl(template, startTag, endTag string, data map[string]any) (string, error) {
	var sb strings.Builder
	_, err := Interpolate(template, startTag, endTag, &sb, func(w io.Writer, tag string) (int, error) {
		if value, ok := data[tag]; ok {
			return write(w, value)
		}
		if value, ok := data["*"]; ok {
			return write(w, value)
		}
		bts := UnsafeStringToBytes(startTag)
		bts = append(bts, tag...)
		bts = append(bts, endTag...)
		return w.Write(bts)
	})
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func write(w io.Writer, value any) (int, error) {
	if value == nil {
		return 0, nil
	}
	switch v := value.(type) {
	case bool:
		return w.(*strings.Builder).WriteString(strconv.FormatBool(v))
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return fmt.Fprintf(w, "%d", v)
	case string:
		return w.(*strings.Builder).WriteString(v)
	case *string:
		return w.(*strings.Builder).WriteString(*v)
	case []byte:
		return w.Write(v)
	case *[]byte:
		return w.Write(*v)
	case fmt.Stringer:
		return w.(*strings.Builder).WriteString(v.String())
	default:
		return fmt.Fprintf(w, "%v", value)
	}
}

// Interpolate calls f on each template tag (placeholder) occurrence.
//
// Returns the number of bytes written to w.
func Interpolate(template, startTag, endTag string, w io.Writer, f TagFunc) (int64, error) {
	s := UnsafeStringToBytes(template)
	a := UnsafeStringToBytes(startTag)
	b := UnsafeStringToBytes(endTag)

	var nn int64
	var ni int
	var err error
	for {
		n := bytes.Index(s, a)
		if n < 0 {
			break
		}
		ni, err = w.Write(s[:n])
		nn += int64(ni)
		if err != nil {
			return nn, err
		}

		s = s[n+len(a):]
		n = bytes.Index(s, b)
		if n < 0 {
			// cannot find end tag - just write it to the output.
			ni, _ = w.Write(a)
			nn += int64(ni)
			break
		}

		ni, err = f(w, UnsafeBytesToString(s[:n]))
		nn += int64(ni)
		if err != nil {
			return nn, err
		}
		s = s[n+len(b):]
	}
	ni, err = w.Write(s)
	nn += int64(ni)

	return nn, err
}
