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

func TestSubstitute_Basic(t *testing.T) {
	out, err := Substitute("Hello {name}, id={id}", map[string]any{
		"name": "Bob",
		"id":   42,
	})
	if err != nil {
		t.Fatalf("Substitute error: %v", err)
	}
	if out != "Hello Bob, id=42" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestSubstitute_StarFallback_AndUnknownKept(t *testing.T) {
	out, err := Substitute("{a}-{b}-{c}", map[string]any{"*": "X", "b": "B"})
	if err != nil {
		t.Fatalf("Substitute error: %v", err)
	}
	// a -> X via * , b -> B, c -> X via *
	if out != "X-B-X" {
		t.Fatalf("unexpected: %q", out)
	}
}

func TestInterpolate_CustomDelimiters_AndTypes(t *testing.T) {
	name := "Ann"
	bytesVal := []byte("BYTES")
	out, err := Interpolate("[[s]]|[[p]]|[[b]]|[[sb]]|[[str]]|[[def]]|[[unk]]", "[[", "]]", map[string]any{
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

func TestTmpl_UnmatchedEndWritesLiteralStart(t *testing.T) {
	var buf strings.Builder
	n, err := Tmpl("pre {no_end tail", "{", "}", &buf, func(w io.Writer, tag string) (int, error) {
		return 0, nil
	})
	if err != nil {
		t.Fatalf("Interpolate error: %v", err)
	}
	if n == 0 || !strings.Contains(buf.String(), "{no_end tail") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestTmpl_CallsTagFunc(t *testing.T) {
	var buf bytes.Buffer
	var seen []string
	_, err := Tmpl("A{one}B{two}C", "{", "}", &buf, func(w io.Writer, tag string) (int, error) {
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

// Benchmark tests moved from bench_test.go
func BenchmarkSubstitute(b *testing.B) {
	in := "Hello, {name}! Today is {day}. {name} likes {thing}."
	m := map[string]any{"name": "Alice", "day": "Wednesday", "thing": "Go"}
	b.ReportAllocs()
	for b.Loop() {
		_, _ = Substitute(in, m)
	}
}

func BenchmarkInterpolate_Large(b *testing.B) {
	// Build a large template with repeated segments
	part := "User:{ID},Name:{Name},Post:{Post},Q={Q};"
	tpl := strings.Repeat(part, 200) // ~200 segments
	data := map[string]any{"ID": 123, "Post": 456, "Q": "search", "Name": "alice"}
	b.ReportAllocs()
	for b.Loop() {
		_, _ = Interpolate(tpl, "{", "}", data)
	}
}

func BenchmarkInterpolate_DensePlaceholders(b *testing.B) {
	// Template with very dense placeholders
	tpl := "{A}{B}{C}{D}{E}{F}{G}{H}{I}{J}{K}{L}{M}{N}{O}{P}{Q}{R}{S}{T}"
	data := map[string]any{}
	for i := 'A'; i <= 'T'; i++ {
		data[string(i)] = int(i)
	}
	b.ReportAllocs()
	for b.Loop() {
		_, _ = Interpolate(tpl, "{", "}", data)
	}
}

func BenchmarkInterpolateTemp(b *testing.B) {
	tpl := "/users/{ID}/posts/{Post}?q={Q}"
	data := map[string]any{"ID": 123, "Post": 456, "Q": "search"}
	b.ReportAllocs()
	for b.Loop() {
		_, _ = Interpolate(tpl, "{", "}", data)
	}
}

func BenchmarkTmpl_TagFunc(b *testing.B) {
	tpl := "Hi, {name}! You have {n} messages."
	vars := map[string]string{"name": "alice", "n": "7"}
	fn := TagFunc(func(w io.Writer, tag string) (int, error) {
		if v, ok := vars[tag]; ok {
			// simulate some work
			return w.Write([]byte(strings.ToUpper(v)))
		}
		return w.Write([]byte(""))
	})
	b.ReportAllocs()
	for b.Loop() {
		_, _ = Tmpl(tpl, "{", "}", &strings.Builder{}, fn)
	}
}
