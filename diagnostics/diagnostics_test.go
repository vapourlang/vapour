package diagnostics

import (
	"testing"

	"github.com/vapourlang/vapour/token"
)

func TestUniqueAcrossFiles(t *testing.T) {
	ds := Diagnostics{
		NewError(token.Item{File: "a.vp", Line: 0, Char: 3}, "a"),
		NewError(token.Item{File: "b.vp", Line: 0, Char: 3}, "b"),
		NewError(token.Item{File: "a.vp", Line: 0, Char: 3}, "dup"),
		// would collide with "%d%d" keys: 1,23 vs 12,3
		NewError(token.Item{File: "a.vp", Line: 1, Char: 23}, "c"),
		NewError(token.Item{File: "a.vp", Line: 12, Char: 3}, "d"),
	}

	if got := len(ds.Unique()); got != 4 {
		t.Fatalf("Unique: expected 4 diagnostics, got %v", got)
	}

	if got := len(ds.UniqueLine()); got != 4 {
		t.Fatalf("UniqueLine: expected 4 diagnostics, got %v", got)
	}
}
