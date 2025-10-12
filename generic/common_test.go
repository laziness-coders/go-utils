package generic

import (
	"reflect"
	"testing"
)

func TestToPointerBool(t *testing.T) {
	v := true
	p := ToPointer(v)
	if p == nil {
		t.Fatalf("expected non-nil pointer")
	}
	if *p != v {
		t.Fatalf("expected pointer to %v, got %v", v, *p)
	}
	if reflect.TypeOf(p) != reflect.TypeOf(&v) {
		t.Fatalf("expected type %T, got %T", &v, p)
	}
}

func TestToPointerInt(t *testing.T) {
	v := 42
	p := ToPointer(v)
	if p == nil {
		t.Fatalf("expected non-nil pointer")
	}
	if *p != v {
		t.Fatalf("expected pointer to %v, got %v", v, *p)
	}
	if reflect.TypeOf(p) != reflect.TypeOf(&v) {
		t.Fatalf("expected type %T, got %T", &v, p)
	}
}

func TestToPointerString(t *testing.T) {
	v := "hello"
	p := ToPointer(v)
	if p == nil {
		t.Fatalf("expected non-nil pointer")
	}
	if *p != v {
		t.Fatalf("expected pointer to %v, got %v", v, *p)
	}
	if reflect.TypeOf(p) != reflect.TypeOf(&v) {
		t.Fatalf("expected type %T, got %T", &v, p)
	}
}

func TestToPointerStruct(t *testing.T) {
	type MyStruct struct {
		A int
		B string
	}
	v := MyStruct{A: 1, B: "test"}
	p := ToPointer(v)
	if p == nil {
		t.Fatalf("expected non-nil pointer")
	}
	if *p != v {
		t.Fatalf("expected pointer to %v, got %v", v, *p)
	}
	if reflect.TypeOf(p) != reflect.TypeOf(&v) {
		t.Fatalf("expected type %T, got %T", &v, p)
	}
}

func TestToPointerSlice(t *testing.T) {
	v := []int{1, 2, 3}
	p := ToPointer(v)
	if p == nil {
		t.Fatalf("expected non-nil pointer")
	}
	if !reflect.DeepEqual(*p, v) {
		t.Fatalf("expected pointer to %v, got %v", v, *p)
	}
	if reflect.TypeOf(p) != reflect.TypeOf(&v) {
		t.Fatalf("expected type %T, got %T", &v, p)
	}
}
