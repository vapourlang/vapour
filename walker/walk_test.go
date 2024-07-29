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

type id: struct {
  int,
  name: char
}

# should fail, already defined
type id: int

# should fail, cannot find type
let z: undefinedType = "hello"

# should fail, different types
let v: int = (10, "hello", na)

# should fail, it's a const
y = 2

# should fail, type mismatch
let wrongType: num = "hello"

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

func TestEnvironment(t *testing.T) {
	code := `
let z: int = 1

func add(n: int, y: int) int | na {
  if n == 1 {
    return na
  }

  return n + y
}

# should fail, this can be na
let result: int = add(1, 2)

const y: int = 1

for(let i: int = 1 in 1:10) {
  print(i)
}
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
