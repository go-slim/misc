package misc

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestMD5(t *testing.T) {
	if MD5("hello") != "5D41402ABC4B2A76B9719D911017C592" {
		t.Fatal("MD5 failed")
	}
	if MD5("") != "D41D8CD98F00B204E9800998ECF8427E" {
		t.Fatal("MD5 empty failed")
	}
}

func TestSha1(t *testing.T) {
	if Sha1("hello") != "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d" {
		t.Fatal("Sha1 failed")
	}
	if Sha1("") != "da39a3ee5e6b4b0d3255bfef95601890afd80709" {
		t.Fatal("Sha1 empty failed")
	}
}

func TestSha256(t *testing.T) {
	if Sha256("hello") != "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824" {
		t.Fatal("Sha256 failed")
	}
	if Sha256("") != "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		t.Fatal("Sha256 empty failed")
	}
}

func TestPasswordHashAndVerify(t *testing.T) {
	pwd := "S3cret!"
	h, err := PasswordHash(pwd)
	if err != nil {
		t.Fatalf("PasswordHash error: %v", err)
	}
	if h == "" {
		t.Fatal("PasswordHash returned empty hash")
	}
	if !PasswordVerify(pwd, h) {
		t.Fatal("PasswordVerify should succeed with correct password")
	}
	if PasswordVerify("wrong", h) {
		t.Fatal("PasswordVerify should fail with wrong password")
	}
}

// Benchmark tests moved from bench_test.go
func BenchmarkMD5(b *testing.B) {
	s := strings.Repeat("x", 1024)
	b.ReportAllocs()
	for b.Loop() {
		_ = MD5(s)
	}
}

func BenchmarkSha1(b *testing.B) {
	s := strings.Repeat("x", 1024)
	b.ReportAllocs()
	for b.Loop() {
		_ = Sha1(s)
	}
}

func BenchmarkSha256(b *testing.B) {
	s := strings.Repeat("x", 1024)
	b.ReportAllocs()
	for b.Loop() {
		_ = Sha256(s)
	}
}

func BenchmarkPasswordHashVerify(b *testing.B) {
	pwd := "S3cret!"
	hash, _ := PasswordHash(pwd)

	b.Run("hash", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = PasswordHash(pwd)
		}
	})

	b.Run("verify", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = PasswordVerify(pwd, hash)
		}
	})
}

func BenchmarkPasswordHash_Cost(b *testing.B) {
	pwd := []byte("S3cret!")
	for _, cost := range []int{8, bcrypt.DefaultCost, 12} {
		b.Run("cost="+strconv.Itoa(cost), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = bcrypt.GenerateFromPassword(pwd, cost)
			}
		})
	}
}

// Example tests for godoc

func ExampleMD5() {
	hash := MD5("hello")
	fmt.Println(hash)
	// Output: 5D41402ABC4B2A76B9719D911017C592
}

func ExampleSha1() {
	hash := Sha1("hello")
	fmt.Println(hash)
	// Output: aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
}

func ExampleSha256() {
	hash := Sha256("hello")
	fmt.Println(hash)
	// Output: 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
}

func ExamplePasswordHash() {
	password := "MySecretPassword123"
	hash, err := PasswordHash(password)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// The hash will be different each time due to salt
	fmt.Println("Hash generated successfully, length:", len(hash))
	// Output: Hash generated successfully, length: 60
}

func ExamplePasswordVerify() {
	password := "MySecretPassword123"
	hash, _ := PasswordHash(password)

	// Verify correct password
	if PasswordVerify(password, hash) {
		fmt.Println("Correct password")
	}

	// Verify wrong password
	if !PasswordVerify("WrongPassword", hash) {
		fmt.Println("Wrong password rejected")
	}
	// Output: Correct password
	// Wrong password rejected
}
