package misc

import "testing"

func TestExtensionByType_Builtin(t *testing.T) {
	ext := ExtensionByType("image/jpeg")
	if ext != ".jpeg" { // first extension in table
		t.Fatalf("ExtensionByType(image/jpeg) = %q, want %q", ext, ".jpeg")
	}
}

func TestTypeByExtension_Builtin(t *testing.T) {
	typ := TypeByExtension(".png")
	if typ != "image/png" {
		t.Fatalf("TypeByExtension(.png) = %q, want %q", typ, "image/png")
	}
}

func TestTypeByExtension_NoDot(t *testing.T) {
	typ := TypeByExtension("jpg")
	if typ != "image/jpeg" {
		t.Fatalf("TypeByExtension(jpg) = %q, want %q", typ, "image/jpeg")
	}
}

func TestExtensionByType_FallbackUnknown(t *testing.T) {
	ext := ExtensionByType("application/x-unknown-type-xyz")
	if ext != "" {
		t.Fatalf("ExtensionByType should return empty for unknown type, got %q", ext)
	}
}

func TestCharsetByType(t *testing.T) {
	if got := CharsetByType("text/plain"); got != "charset=utf-8" {
		t.Fatalf("CharsetByType(text/plain) = %q, want charset=utf-8", got)
	}
	if got := CharsetByType("application/json"); got != "" {
		t.Fatalf("CharsetByType(application/json) = %q, want empty", got)
	}
}
