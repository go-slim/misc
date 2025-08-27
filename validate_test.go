package misc

import "testing"

type myStruct struct{ A int }

func TestIsZero_Basics(t *testing.T) {
	if !IsZero(0) || !IsZero(0.0) || !IsZero("") || !IsZero(false) {
		t.Fatal("basic zero checks failed")
	}
	if IsZero(1) || IsZero("x") || IsZero(true) {
		t.Fatal("basic non-zero checks failed")
	}
}

func TestIsZero_PointersAndStructs(t *testing.T) {
	var p *int
	if !IsZero(p) {
		t.Fatal("nil pointer should be zero")
	}
	n := 0
	pp := &p
	if !IsZero(pp) { // pointer to nil pointer -> zero
		t.Fatal("pointer to nil pointer should be zero")
	}
	p = &n
	if !IsZero(p) { // pointer to zero value int -> still zero
		t.Fatal("pointer to zero value should be zero")
	}
	m := myStruct{}
	if !IsZero(m) || IsZero(myStruct{A: 1}) {
		t.Fatal("struct zero checks failed")
	}
}
