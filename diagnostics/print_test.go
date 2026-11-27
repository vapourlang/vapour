package diagnostics

import (
	"strings"
	"testing"

	"github.com/vapourlang/vapour/token"
)

func TestStringPosition(t *testing.T) {
	// positions are 0-based and Char is the end of the token
	d := NewWarning(token.Item{File: "a.vp", Line: 3, Char: 10, Value: "nope"}, "not found")

	if !strings.Contains(d.String(), "a.vp:4:7 ") {
		t.Fatalf("expected 1-based line and start column a.vp:4:7, got %v", d.String())
	}
}
