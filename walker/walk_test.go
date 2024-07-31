package walker

import (
	"fmt"
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

	fmt.Println("-----------------------------")
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

	fmt.Println("-----------------------------")
	w.Walk(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestInfix(t *testing.T) {
	code := `let x: char = "hello"

# should fail, cannot be NA
x = na

# should fail, types do not match
let z: char = 1

func add(n: int, y: int) int | na {
	if n == 1 {
		return na
	}

  return n + y
}

# should fail, this can be na
let result: int = add(1, 2)

# should fail, const must have single type
const v: int | na = 1
`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("-----------------------------")
	w.Walk(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestFunction(t *testing.T) {
	code := `
func foo(n: int) int {
  let x: char = "hello"

  # should fail, returns wrong type
  if (n == 2) {
    return "hello"
  }

  if (n == 3) {
    return 1.2
  }

  # should fail, returns wrong type
  return x
}
`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("-----------------------------")
	w.Walk(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestNumber(t *testing.T) {
	code := `let x: num = 1
x = 1.1

let u: int = 1e10

let integer: int = 1;

# should fail, assign flaot to int
integer = 2.1
`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("-----------------------------")
	w.Walk(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestAnonymous(t *testing.T) {
	code := `
# should fail, returns wrong type
lapply(1..10, (z: int): int => {
  return "hello"
}) `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("-----------------------------")
	w.Walk(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}
