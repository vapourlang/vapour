package parser

import (
	"fmt"
	"testing"

	"github.com/devOpifex/vapour/lexer"
)

func TestBasic(t *testing.T) {
	code := `let x = 1  `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	l.Print()
	p := New(l)

	prog := p.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where parsed")
	}

	fmt.Println("+++++++++")
	fmt.Printf("number of statements: %v\n", len(prog.Statements))
	fmt.Println(prog.String())
}
