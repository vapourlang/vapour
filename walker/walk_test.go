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

# should fail, must have a value
const xx: int
`

	l := lexer.NewTest(code)

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

func addz(n: int, y: int): int | na {
	if n == 1 {
		return na
	}

	return n + y
}

# should fail, this can be na
let result: int = addz(1, 2)

const y: int = 1

for(let i: int = 1 in 1:10) {
  print(i)
}
`

	l := lexer.NewTest(code)

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

v = 2
`

	l := lexer.NewTest(code)

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

	l := lexer.NewTest(code)

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

	l := lexer.NewTest(code)

	l.Run()
	fmt.Println("-----------------------------")
	p := parser.New(l)

	prog := p.Run()

	w := New()
	w.Walk(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestNamespace(t *testing.T) {
	code := `# should fail, duplicated params
func bar(x: int, x: int): int {return x + y}
`

	l := lexer.NewTest(code)

	l.Run()
	fmt.Println("-----------------------------")
	p := parser.New(l)

	prog := p.Run()

	w := New()
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

let integer: int = 1

# should fail, assign num to int
integer = 2.1

let x: int = sum(1,2,3)
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	w.Walk(prog)

	fmt.Println("-----------------------------")
	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestMethod(t *testing.T) {
	code := `func (o: obj) add(n: int): char {
  return "hello"
}

type person: struct{
  int,
	name: char
}

# should fail, xxx does not exist
func (p: person) setName(name: char): null {
  p$xxx = 2
}

# should fail, name expects char
func (p: person) setName(name: char): null {
  p$name = 2
}
`

	l := lexer.NewTest(code)

	l.Run()
	fmt.Println("-----------------------------")
	p := parser.New(l)

	prog := p.Run()

	w := New()
	w.Walk(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestCall(t *testing.T) {
	code := `func foo(x: int, y: char): int {
  print(y)
	return x + 1
}

# should fail, argument does not exist
foo(z = 2)

# should fail, wrong type
foo(
  x = "hello"
)

# should fail, wrong type
foo("hello")

# should fail, too many arguments
foo(1, "hello", 3)

func lg (...: char): null {
  print(...)
}

lg("hello", "world")

# should fail, wrong type
lg("hello", 1)
`

	l := lexer.NewTest(code)

	l.Run()
	fmt.Println("-----------------------------")
	p := parser.New(l)

	prog := p.Run()

	w := New()
	w.Walk(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}
