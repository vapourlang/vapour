package walker

import (
	"testing"

	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
)

func TestBasic(t *testing.T) {
	code := `let x: int | na = 1

x = 2

# should fail, it's already declared
let x: char = "hello"

const y: int = 1

# should fail, it's a const
y = 2

type id: struct {
  int,
  name: char
}

# should fail, already defined
type id: int

# should fail, cannot find type
let z: undefinedType = "hello"

# should fail, type mismatch
let wrongType: num = "hello"

func add(x: int, y: int) int | na {
  if x == 1 {
    return na
  }

  return x + y
}

# should fail, this can be na
let result: int = add(1, 2)
`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	w.Walk(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}
