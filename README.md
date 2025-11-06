# misc

![CI](https://github.com/go-slim/misc/actions/workflows/ci.yml/badge.svg)

A general-purpose utility library providing common tools and types, covering digest/password, functional composition, MIME parsing, template interpolation, stack information, zero-copy conversion, and zero value checks.

[中文文档](README.zh-CN.md)

- Module path: `go-slim.dev/misc`
- Go version: `1.24`

## Features

- **Digest/Password**: `MD5`, `Sha1`, `Sha256`, `PasswordHash`/`PasswordVerify`
- **Function composition**: `Call`/`CallG`, `Wrap`/`WrapG`
- **MIME parsing**: `ExtensionByType`, `TypeByExtension`, `CharsetByType` **[Deprecated]**
- **Template interpolation**: `Substitute`/`Interpolate`, `Tmpl`, `TagFunc`
- **Generic utilities**: `Zero`, `Ptr`, `Nil`, `IsZero`, `IsNil`, `Coalesce`
- **Math utilities**: `MinMax`, `Clamp`
- **Zero-copy conversion**: `BytesToString`, `StringToBytes` (use with caution)
- **Stack information**: `Stack()` (includes source code lines) **[Deprecated]**

## Installation

```bash
go get go-slim.dev/misc
```

## Quick Examples

### Digest and Password

```go
package main

import (
    "fmt"
    "go-slim.dev/misc"
)

func main() {
    fmt.Println(misc.MD5("hello"))
    fmt.Println(misc.Sha1("hello"))
    fmt.Println(misc.Sha256("hello"))

    hash := misc.PasswordHash("S3cret!")
    ok := misc.PasswordVerify("S3cret!", hash)
    fmt.Println("password ok:", ok)
}
```

### Function Composition

```go
// Call executes functions sequentially, stops on first error
err := misc.Call(
    func() error { fmt.Println("step 1"); return nil },
    func() error { fmt.Println("step 2"); return nil },
)

// CallG calls multiple functions with the same parameter
err = misc.CallG(42,
    func(v int) error { fmt.Println("value:", v); return nil },
    func(v int) error { fmt.Println("doubled:", v*2); return nil },
)

// Wrap combines multiple functions into one
wrapped := misc.Wrap(
    func() error { return nil },
    func() error { return nil },
)
_ = wrapped()
_ = err
```

### MIME Parsing (Deprecated)

```go
ext := misc.ExtensionByType("image/png")      // .png
typ := misc.TypeByExtension(".json")          // application/json
cs  := misc.CharsetByType("text/plain")       // charset=utf-8
```

### Template Interpolation

```go
// Substitute uses {key} format for replacement
s, _ := misc.Substitute("Hello, {name}!", map[string]any{"name": "Alice"})
fmt.Println(s) // Hello, Alice!

// Interpolate uses custom tags with map data
out, _ := misc.Interpolate("/user/{{ID}}?q={{Q}}", "{{", "}}", map[string]any{"ID": 7, "Q": "k"})
fmt.Println(out) // /user/7?q=k

// Tmpl is the low-level function supporting custom TagFunc
var buf bytes.Buffer
misc.Tmpl("value: {{x}}", "{{", "}}", &buf, func(w io.Writer, tag string) (int, error) {
    return w.Write([]byte("42"))
})
fmt.Println(buf.String()) // value: 42
```

### Generic Utilities

```go
// Zero value
zeroInt := misc.Zero[int]()        // 0
zeroStr := misc.Zero[string]()     // ""

// Pointer creation
ptr := misc.Ptr(42)                // *int pointing to 42
strPtr := misc.Ptr("hello")        // *string pointing to "hello"

// Nil pointer
var intPtr *int = misc.Nil[int]()  // nil

// Zero check
misc.IsZero(0)                     // true
misc.IsZero("")                    // true
misc.IsZero((*int)(nil))           // true
misc.IsZero(42)                    // false

// Nil check
misc.IsNil(nil)                    // true
misc.IsNil((*int)(nil))            // true
misc.IsNil(42)                     // false

// Coalesce (first non-zero value)
result := misc.Coalesce("", "default", "fallback") // "default"
result := misc.Coalesce(0, 42, 100)                // 42
```

### Math Utilities

```go
// MinMax returns (min, max)
min, max := misc.MinMax(5, 3)           // (3, 5)
min, max := misc.MinMax(1.5, 2.7)       // (1.5, 2.7)
min, max := misc.MinMax("b", "a")       // ("a", "b")

// Clamp constrains value to range
result := misc.Clamp(15, 10, 20)        // 15 (within range)
result := misc.Clamp(5, 10, 20)         // 10 (below minimum)
result := misc.Clamp(25, 10, 20)        // 20 (above maximum)
result := misc.Clamp(15, 20, 10)        // 15 (auto-swaps bounds)
```

### Zero-Copy Conversion (Unsafe)

```go
// Warning: Use only when you fully understand the risks!
b := []byte("hello")
s := misc.BytesToString(b)       // Zero-copy conversion
bs := misc.StringToBytes(s)      // Zero-copy conversion

// WARNING: Do not modify the original data after conversion!
// They share the same underlying memory.
```

### Stack Trace (Deprecated)

```go
trace := misc.Stack(0) // Skip 0 frames, includes function, file, and line number
fmt.Println(trace)
// Use runtime.Stack or debug.Stack instead
```

## Running Tests

```bash
go test ./...
```

## License

This project is licensed under the MIT License. See `LICENSE` for details.
