package misc

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

type testStringer int

func (t testStringer) String() string { return fmt.Sprintf("S%d", int(t)) }

func TestStrtr_Basic(t *testing.T) {
	out, err := Strtr("Hello {name}, id={id}", map[string]any{
		"name": "Bob",
		"id":   42,
	})
	if err != nil {
		t.Fatalf("Strtr error: %v", err)
	}
	if out != "Hello Bob, id=42" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestStrtr_StarFallback_AndUnknownKept(t *testing.T) {
	out, err := Strtr("{a}-{b}-{c}", map[string]any{"*": "X", "b": "B"})
	if err != nil {
		t.Fatalf("Strtr error: %v", err)
	}
	// a -> X via * , b -> B, c -> X via *
	if out != "X-B-X" {
		t.Fatalf("unexpected: %q", out)
	}
}

func TestTmpl_CustomDelimiters_AndTypes(t *testing.T) {
	name := "Ann"
	bytesVal := []byte("BYTES")
	out, err := Tmpl("[[s]]|[[p]]|[[b]]|[[sb]]|[[str]]|[[def]]|[[unk]]", "[[", "]]", map[string]any{
		"s":   "hi",
		"p":   &name,
		"b":   bytesVal,
		"sb":  &bytesVal,
		"str": testStringer(7),
		"def": struct{ X int }{3},
	})
	if err != nil {
		t.Fatalf("Tmpl error: %v", err)
	}
	want := "hi|Ann|BYTES|BYTES|S7|{3}|[[unk]]"
	if out != want {
		t.Fatalf("got %q want %q", out, want)
	}
}

func TestInterpolate_UnmatchedEndWritesLiteralStart(t *testing.T) {
	var buf strings.Builder
	n, err := Interpolate("pre {no_end tail", "{", "}", &buf, func(w io.Writer, tag string) (int, error) {
		return 0, nil
	})
	if err != nil {
		t.Fatalf("Interpolate error: %v", err)
	}
	if n == 0 || !strings.Contains(buf.String(), "{no_end tail") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestInterpolate_CallsTagFunc(t *testing.T) {
	var buf bytes.Buffer
	var seen []string
	_, err := Interpolate("A{one}B{two}C", "{", "}", &buf, func(w io.Writer, tag string) (int, error) {
		seen = append(seen, tag)
		return w.Write([]byte(strings.ToUpper(tag)))
	})
	if err != nil {
		t.Fatalf("Interpolate error: %v", err)
	}
	if strings.Join(seen, ",") != "one,two" {
		t.Fatalf("tags seen wrong: %v", seen)
	}
	if buf.String() != "AONEBTWOC" {
		t.Fatalf("output wrong: %q", buf.String())
	}
}
