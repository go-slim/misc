# misc

![CI](https://github.com/go-slim/misc/actions/workflows/ci.yml/badge.svg)

通用实用库，提供常用的工具函数与类型，覆盖摘要/密码、函数式封装、MIME 解析、模板插值、堆栈信息、零拷贝转换与零值判断等能力。

- 模块路径：`go-slim.dev/misc`
- Go 版本：`1.24`

## 功能概览

- 摘要/密码：`MD5`、`Sha1`、`Sha256`、`PasswordHash`/`PasswordVerify`
- 函数组合：`Call`/`CallG`、`Wrap`/`WrapG`
- MIME 解析：`ExtensionByType`、`TypeByExtension`、`CharsetByType` **[已废弃]**
- 模板插值：`Substitute`/`Interpolate`、`Tmpl`、`TagFunc`
- 泛型工具：`Zero`、`Ptr`、`Nil`、`IsZero`、`IsNil`、`Coalesce`
- 数学工具：`MinMax`、`Clamp`
- 零拷贝转换：`BytesToString`、`StringToBytes`（需谨慎使用）
- 堆栈信息：`Stack()`（包含源码行）**[已废弃]**

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

### MIME 解析（已废弃）

```go
ext := misc.ExtensionByType("image/png")      // .png
typ := misc.TypeByExtension(".json")          // application/json
cs  := misc.CharsetByType("text/plain")       // charset=utf-8
```

### 模板插值

```go
// Substitute 使用 {key} 格式替换
s, _ := misc.Substitute("Hello, {name}!", map[string]any{"name": "Alice"})
fmt.Println(s) // Hello, Alice!

// Interpolate 使用自定义标签和 map 数据
out, _ := misc.Interpolate("/user/{{ID}}?q={{Q}}", "{{", "}}", map[string]any{"ID": 7, "Q": "k"})
fmt.Println(out) // /user/7?q=k

// Tmpl 是底层函数，支持自定义 TagFunc
var buf bytes.Buffer
misc.Tmpl("value: {{x}}", "{{", "}}", &buf, func(w io.Writer, tag string) (int, error) {
    return w.Write([]byte("42"))
})
fmt.Println(buf.String()) // value: 42
```

### 泛型工具

```go
// 零值
zeroInt := misc.Zero[int]()        // 0
zeroStr := misc.Zero[string]()     // ""

// 创建指针
ptr := misc.Ptr(42)                // *int 指向 42
strPtr := misc.Ptr("hello")        // *string 指向 "hello"

// Nil 指针
var intPtr *int = misc.Nil[int]()  // nil

// 零值检查
misc.IsZero(0)                     // true
misc.IsZero("")                    // true
misc.IsZero((*int)(nil))           // true
misc.IsZero(42)                    // false

// Nil 检查
misc.IsNil(nil)                    // true
misc.IsNil((*int)(nil))            // true
misc.IsNil(42)                     // false

// Coalesce（返回第一个非零值）
result := misc.Coalesce("", "default", "fallback") // "default"
result := misc.Coalesce(0, 42, 100)                // 42
```

### 数学工具

```go
// MinMax 返回 (最小值, 最大值)
min, max := misc.MinMax(5, 3)           // (3, 5)
min, max := misc.MinMax(1.5, 2.7)       // (1.5, 2.7)
min, max := misc.MinMax("b", "a")       // ("a", "b")

// Clamp 将值限制在范围内
result := misc.Clamp(15, 10, 20)        // 15（在范围内）
result := misc.Clamp(5, 10, 20)         // 10（低于最小值）
result := misc.Clamp(25, 10, 20)        // 20（高于最大值）
result := misc.Clamp(15, 20, 10)        // 15（自动交换边界）
```

### 零拷贝转换（危险）

```go
// 注意：仅在明确理解风险时使用！
b := []byte("hello")
s := misc.BytesToString(b)       // 零拷贝转换
bs := misc.StringToBytes(s)      // 零拷贝转换

// 警告：转换后不要修改原始数据！
// 它们共享相同的底层内存。
```

### 堆栈追踪（已废弃）

```go
trace := misc.Stack(0) // 跳过 0 层，包含函数、文件与行号
fmt.Println(trace)
// 建议使用 runtime.Stack 或 debug.Stack 替代
```

## 运行测试

```bash
go test ./...
```

## 许可证

本项目使用 MIT License，详见 `LICENSE`。
