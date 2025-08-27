package misc

import (
	"strings"
	"testing"
)

func stackHelper() string { return Stack(0) }

func TestStackBasic(t *testing.T) {
	s := stackHelper()
	if s == "" {
		t.Fatal("Stack returned empty string")
	}
	// Should include current file or helper function name line
	if !strings.Contains(s, "stack_test.go") && !strings.Contains(s, "stackHelper") {
		t.Fatalf("stack output seems unexpected: %q", s)
	}
}
