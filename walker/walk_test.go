package walker

import (
	"testing"

	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
)

func TestBasic(t *testing.T) {
	code := `let x: int | num = 1  `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	w.Walk(prog)
}
