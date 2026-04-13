package renderer

import (
	"strings"
	"testing"
)

func TestFormatGo(t *testing.T) {
	t.Parallel()

	t.Run("formats valid Go code", func(t *testing.T) {
		t.Parallel()

		src := []byte(`package main

import "fmt"

func main() {
fmt.Println(  "hello"  )
}
`)

		result, err := FormatGo(src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.Contains(string(result), `fmt.Println("hello")`) {
			t.Errorf("expected formatted output, got:\n%s", result)
		}
	})

	t.Run("returns error for invalid Go code", func(t *testing.T) {
		t.Parallel()

		src := []byte(`package main

func {{{ invalid
`)

		_, err := FormatGo(src)
		if err == nil {
			t.Error("expected error for invalid Go code")
		}
	})
}

func TestIndicateLines(t *testing.T) {
	t.Parallel()

	src := []byte("line one\nline two\nline three")
	result := indicateLines(src)

	if !strings.Contains(result, "   1: line one") {
		t.Errorf("expected line numbers, got:\n%s", result)
	}
	if !strings.Contains(result, "   3: line three") {
		t.Errorf("expected line 3, got:\n%s", result)
	}
}
