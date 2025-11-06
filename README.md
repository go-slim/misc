# misc

![CI](https://github.com/go-slim/misc/actions/workflows/ci.yml/badge.svg)

通用实用库，提供常用的工具函数与类型，覆盖摘要/密码、函数式封装、MIME 解析、模板插值、堆栈信息、零拷贝转换与零值判断等能力。

- 模块路径：`go-slim.dev/misc`
- Go 版本：`1.24`

## 功能概览

- 摘要/密码：`MD5`、`Sha1`、`Sha256`、`PasswordHash`/`PasswordVerify`
- 函数组合：`Call`/`CallG`、`Wrap`/`WrapG`
- MIME 解析：`ExtensionByType`、`TypeByExtension`、`CharsetByType`
- 模板插值：`Strtr`/`Tmpl`、`Interpolate`、`TagFunc`
- 堆栈信息：`Stack()`（包含源码行）**[已废弃]**
- 零拷贝转换：`UnsafeBytesToString`、`UnsafeStringToBytes`（需谨慎使用）
- 零值判断：`IsZero`（支持指针递归判断）

## 安装

```bash
go get go-slim.dev/misc
```

## 快速示例

### 摘要与密码

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

### 函数组合/封装

```go
v, err := misc.Call(func() (int, error) { return 42, nil })
_ = misc.MustCall(func() (string, error) { return "ok", nil })

wrapped := misc.Wrap(func(a int) int { return a * 2 })
result := wrapped(21) // 42
_ = result
```

### MIME 解析

```go
ext := misc.ExtensionByType("application/json; charset=utf-8") // .json
typ := misc.TypeByExtension(".png")                            // image/png
cs  := misc.CharsetByType("text/html; charset=utf-8")          // utf-8
_ = []string{ext, typ, cs}
```

### 模板插值

```go
s := misc.Strtr("Hello, {name}!", map[string]string{"{name}": "Alice"})

out, _ := misc.Tmpl("/user/{{.ID}}?q={{.Q}}", map[string]any{"ID": 7, "Q": "k"})

// 高级：自定义标签函数
upper := misc.TagFunc("upper", func(s string) string { return strings.ToUpper(s) })
res, _ := misc.Interpolate("Hi, {{upper .Name}}", map[string]any{"Name": "bob"}, upper)
_ = []string{s, out, res}
```

### 堆栈追踪

```go
trace := misc.Stack(0) // 跳过 0 层，包含函数、文件与行号
fmt.Println(trace)
```

### 零拷贝转换（危险）

```go
// 注意：仅在明确理解风险时使用！
b := []byte("hello")
s := misc.UnsafeBytesToString(b)
bs := misc.UnsafeStringToBytes(s)
_ = []any{s, bs}
```

### 零值判断

```go
type T struct{ N int }
var p *T
_ = misc.IsZero(0)   // true
_ = misc.IsZero("")  // true
_ = misc.IsZero(p)   // true（nil 指针）
_ = misc.IsZero(T{}) // true
```

## 运行测试

```bash
go test ./...
```

## 许可证

本项目使用 MIT License，详见 `LICENSE`。
