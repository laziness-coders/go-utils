package myerrors

import (
	"errors"
	"regexp"
	"testing"
)

func TestNewIncludesStack(t *testing.T) {
	err := New("boom")
	if err == nil {
		t.Fatalf("expected error")
	}
	pattern := regexp.MustCompile(`errors_test.go:\d+`)
	if pattern.MatchString(err.Error()) {
		t.Fatalf("expected Error() without stack trace, got %s", err.Error())
	}
	e, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error type")
	}
	if !pattern.MatchString(e.ErrorWithFrame()) {
		t.Fatalf("expected stack trace, got %s", e.ErrorWithFrame())
	}
}

func TestWrap(t *testing.T) {
	orig := errors.New("orig")
	err := Wrap(orig, "wrap")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, orig) {
		t.Fatalf("expected wrapped error to contain original")
	}
	pattern := regexp.MustCompile(`errors_test.go:\d+`)
	if pattern.MatchString(err.Error()) {
		t.Fatalf("expected Error() without stack trace, got %s", err.Error())
	}
	e, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error type")
	}
	if !pattern.MatchString(e.ErrorWithFrame()) {
		t.Fatalf("expected stack trace, got %s", e.ErrorWithFrame())
	}
}

func TestErrorfIncludesStack(t *testing.T) {
	err := Errorf("number %d", 7)
	if err == nil {
		t.Fatalf("expected error")
	}
	pattern := regexp.MustCompile(`errors_test.go:\d+`)
	if pattern.MatchString(err.Error()) {
		t.Fatalf("expected Error() without stack trace, got %s", err.Error())
	}
	var e *Error
	ok := errors.As(err, &e)
	if !ok {
		t.Fatalf("expected *Error type")
	}
	if !pattern.MatchString(e.ErrorWithFrame()) {
		t.Fatalf("expected stack trace, got %s", e.ErrorWithFrame())
	}
}
