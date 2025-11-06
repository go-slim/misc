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
// Call 按顺序执行多个函数，遇到错误立即返回
err := misc.Call(
    func() error { fmt.Println("step 1"); return nil },
    func() error { fmt.Println("step 2"); return nil },
)

// CallG 使用相同参数调用多个函数
err = misc.CallG(42,
    func(v int) error { fmt.Println("value:", v); return nil },
    func(v int) error { fmt.Println("doubled:", v*2); return nil },
)

// Wrap 将多个函数组合成一个函数
wrapped := misc.Wrap(
    func() error { return nil },
    func() error { return nil },
)
_ = wrapped()
_ = err
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
// Strtr 使用 {key} 格式替换
s, _ := misc.Strtr("Hello, {name}!", map[string]any{"name": "Alice"})
fmt.Println(s) // Hello, Alice!

// Tmpl 使用自定义标签
out, _ := misc.Tmpl("/user/{{ID}}?q={{Q}}", "{{", "}}", map[string]any{"ID": 7, "Q": "k"})
fmt.Println(out) // /user/7?q=k

// Interpolate 支持自定义 TagFunc
var buf bytes.Buffer
misc.Interpolate("value: {{x}}", "{{", "}}", &buf, func(w io.Writer, tag string) (int, error) {
    return w.Write([]byte("42"))
})
fmt.Println(buf.String()) // value: 42
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
