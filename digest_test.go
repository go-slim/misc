package misc

import (
	"strings"
	"testing"
)

func TestMD5(t *testing.T) {
	got := MD5("abc")
	// echo -n abc | md5 => 900150983cd24fb0d6963f7d28e17f72 (uppercase expected)
	want := strings.ToUpper("900150983cd24fb0d6963f7d28e17f72")
	if got != want {
		t.Fatalf("MD5(abc) = %s, want %s", got, want)
	}
}

func TestSha1(t *testing.T) {
	got := Sha1("abc")
	want := "a9993e364706816aba3e25717850c26c9cd0d89d"
	if got != want {
		t.Fatalf("Sha1(abc) = %s, want %s", got, want)
	}
}

func TestSha256(t *testing.T) {
	got := Sha256("abc")
	want := "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	if got != want {
		t.Fatalf("Sha256(abc) = %s, want %s", got, want)
	}
}

func TestPasswordHashAndVerify(t *testing.T) {
	pwd := "s3cr3t!"
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
