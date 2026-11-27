package lexer

import (
	"testing"

	"github.com/vapourlang/vapour/token"
)

func TestNumberLiterals(t *testing.T) {
	cases := []struct {
		code  string
		class token.ItemType
	}{
		{"1", token.ItemInteger},
		{"1L", token.ItemInteger},
		{"0x10", token.ItemInteger},
		{"0xFFL", token.ItemInteger},
		{"1e3", token.ItemInteger},
		{"1.5", token.ItemFloat},
		{"1e-3", token.ItemFloat},
		{"2.5e2", token.ItemFloat},
		{"1i", token.ItemFloat},
		{"2.5i", token.ItemFloat},
	}

	for _, c := range cases {
		l := NewTest(c.code + "\n")
		l.Run()

		if len(l.Items) == 0 {
			t.Fatalf("%v: no items lexed", c.code)
		}

		got := l.Items[0]

		if got.Class != c.class || got.Value != c.code {
			t.Fatalf("%v: expected a single %v token, got %v %q", c.code, c.class, got.Class, got.Value)
		}
	}
}

func TestIdentifierRange(t *testing.T) {
	l := NewTest("let x: int = n..10\nlet y: int = my.var..n\n")
	l.Run()

	var values []string
	for _, item := range l.Items {
		if item.Class == token.ItemIdent || item.Class == token.ItemRange {
			values = append(values, item.Value)
		}
	}

	expected := []string{"x", "n", "..", "y", "my.var", "..", "n"}

	if len(values) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, values)
	}

	for i := range expected {
		if values[i] != expected[i] {
			t.Fatalf("expected %v, got %v", expected, values)
		}
	}
}
