package misc

import (
    "io"
    "strconv"
    "strings"
    "testing"

    "golang.org/x/crypto/bcrypt"
)

func BenchmarkStrtr(b *testing.B) {
    in := "Hello, {name}! Today is {day}. {name} likes {thing}."
    m := map[string]any{"name": "Alice", "day": "Wednesday", "thing": "Go"}
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        _, _ = Strtr(in, m)
    }
}

func BenchmarkTmpl_Large(b *testing.B) {
    // Build a large template with repeated segments
    part := "User:{ID},Name:{Name},Post:{Post},Q={Q};"
    tpl := strings.Repeat(part, 200) // ~200 segments
    data := map[string]any{"ID": 123, "Post": 456, "Q": "search", "Name": "alice"}
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        _, _ = Tmpl(tpl, "{", "}", data)
    }
}

func BenchmarkTmpl_DensePlaceholders(b *testing.B) {
    // Template with very dense placeholders
    tpl := "{A}{B}{C}{D}{E}{F}{G}{H}{I}{J}{K}{L}{M}{N}{O}{P}{Q}{R}{S}{T}"
    data := map[string]any{}
    for i := 'A'; i <= 'T'; i++ {
        data[string(i)] = int(i)
    }
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        _, _ = Tmpl(tpl, "{", "}", data)
    }
}

func BenchmarkTmpl(b *testing.B) {
    tpl := "/users/{ID}/posts/{Post}?q={Q}"
    data := map[string]any{"ID": 123, "Post": 456, "Q": "search"}
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        _, _ = Tmpl(tpl, "{", "}", data)
    }
}

func BenchmarkInterpolate_TagFunc(b *testing.B) {
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
    for i := 0; i < b.N; i++ {
        _, _ = Interpolate(tpl, "{", "}", &strings.Builder{}, fn)
    }
}

func BenchmarkMD5(b *testing.B) {
    s := strings.Repeat("x", 1024)
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        _ = MD5(s)
    }
}

func BenchmarkSha1(b *testing.B) {
    s := strings.Repeat("x", 1024)
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        _ = Sha1(s)
    }
}

func BenchmarkSha256(b *testing.B) {
    s := strings.Repeat("x", 1024)
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        _ = Sha256(s)
    }
}

func BenchmarkPasswordHashVerify(b *testing.B) {
    pwd := "S3cret!"
    hash, _ := PasswordHash(pwd)

    b.Run("hash", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            _, _ = PasswordHash(pwd)
        }
    })

    b.Run("verify", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            _ = PasswordVerify(pwd, hash)
        }
    })
}

func BenchmarkUnsafeConversions(b *testing.B) {
    b.Run("bytes_to_string", func(b *testing.B) {
        bs := []byte(strings.Repeat("a", 64))
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            _ = UnsafeBytesToString(bs)
        }
    })

    b.Run("string_to_bytes", func(b *testing.B) {
        s := strings.Repeat("a", 64)
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            _ = UnsafeStringToBytes(s)
        }
    })
}

func BenchmarkIsZero(b *testing.B) {
    type T struct{ N int }
    cases := []any{0, "", (*T)(nil), T{}, []int(nil)}
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        for _, v := range cases {
            _ = IsZero(v)
        }
    }
}

func BenchmarkPasswordHash_Cost(b *testing.B) {
    pwd := []byte("S3cret!")
    for _, cost := range []int{8, bcrypt.DefaultCost, 12} {
        b.Run("cost="+strconv.Itoa(cost), func(b *testing.B) {
            b.ReportAllocs()
            for i := 0; i < b.N; i++ {
                _, _ = bcrypt.GenerateFromPassword(pwd, cost)
            }
        })
    }
}
