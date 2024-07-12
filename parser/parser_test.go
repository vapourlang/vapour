package parser

import (
	"fmt"
	"testing"

	"github.com/devOpifex/vapour/lexer"
)

func TestBasic(t *testing.T) {
	code := `let x: int = 1  `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where parsed")
	}

	fmt.Println(prog.String())
}

func TestFunc(t *testing.T) {
	code := `func add(x: int = 1, y: int = 2) int {
  let total: int = x + y
  return total
}`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where parsed")
	}

	fmt.Println(prog.String())
}
