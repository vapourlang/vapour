package parser

import (
	"fmt"
	"testing"

	"github.com/devOpifex/vapour/lexer"
)

func TestBasic(t *testing.T) {
	code := `let x = 1 `

	l := &lexer.Lexer{
		Input: code,
	}

	p := New(l)

	p.l.Run()
	p.l.Print()

	prog := p.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where parsed")
	}

	fmt.Println("+++++++++")
	fmt.Printf("statements %v", len(prog.Statements))
	fmt.Println(prog.String())
}
