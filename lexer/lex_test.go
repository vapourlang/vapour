package lexer

import (
	"testing"

	"github.com/devOpifex/vapour/token"
)

func TestTypes(t *testing.T) {
	code := `const x: int | na = 1`

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	l.Print()

	expectLength := 8
	if len(l.Items) != expectLength {
		t.Fatalf("Expecting %v tokens, go %v", expectLength, len(l.Items))
	}

	if l.getItem(0).Class != token.ItemConst {
		t.Fatalf("Expecting constant, got %v - %v", l.getItem(0).String(), l.getItem(0).Value)
	}

	if l.getItem(5).Class != token.ItemTypes {
		t.Fatalf("Expecting type, got %v - %v", l.getItem(5).String(), l.getItem(5).Value)
	}
}
