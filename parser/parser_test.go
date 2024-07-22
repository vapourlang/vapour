package parser

import (
	"fmt"
	"testing"

	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/lexer"
)

func TestBasic(t *testing.T) {
	code := `let x: int | num = 1  `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	env := ast.NewEnvironment()
	fmt.Println(prog.Transpile(env))
}

func TestFunc(t *testing.T) {
	code := `func add(x: int = 1, y: int = 2) int {
  let total: int = x + y * 2
  return total
}`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	env := ast.NewEnvironment()
	fmt.Println(prog.Transpile(env))
}

func TestPipe(t *testing.T) {
	code := `func add() {
  df |>
    mutate(x = 1)
}`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	env := ast.NewEnvironment()
	fmt.Println(prog.Transpile(env))
}

func TestString(t *testing.T) {
	code := `let x: string <- "a \"string\""
let y: string <- 'single quotes'`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	env := ast.NewEnvironment()
	fmt.Println(prog.Transpile(env))
}

func TestComment(t *testing.T) {
	code := `#' @return something
func add() int | number {
  # compute stuff
  let x: tibble = df |>
    mutate(
      x = "hello",
      y = na,
      b = true
    ) |>
    select(x)

  return x
}`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	env := ast.NewEnvironment()
	fmt.Println(prog.Transpile(env))
}

func TestMethod(t *testing.T) {
	code := `func (o obj) add(n: int) string {
  return "hello"
}`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	env := ast.NewEnvironment()
	fmt.Println(prog.Transpile(env))
}

func TestTypeDeclaration(t *testing.T) {
	code := `type userId: int

type obj: struct {
  int | string,
  name: string,
  id: int
} `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	env := ast.NewEnvironment()
	fmt.Println(prog.Transpile(env))
}

func TestAnonymous(t *testing.T) {
	code := `const x: string = "world"
lapply(("hello", x), (x: string) null => { print(x)}) `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	env := ast.NewEnvironment()
	fmt.Println(prog.Transpile(env))
}

func TestIdent(t *testing.T) {
	code := `let x: int = (1,2,3)

x[1, 2] = 15

x[[3]] = 15

df.x = 23

print(x) `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	env := ast.NewEnvironment()
	fmt.Println(prog.Transpile(env))
}

func TestTElipsis(t *testing.T) {
	code := `func foo(...: any) char {
  paste0(..., collapse = ", ")
}`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	env := ast.NewEnvironment()
	fmt.Println(prog.Transpile(env))
}

func TestS3(t *testing.T) {
	code := `
type person: struct {
  int | num,
  name: string,
  age: int
}

func (p person) getAge() int {
  return p.age
}

func (p person) setAge(n: int) null {
  p.age = n
}

func create(name: string, age: int) person {
  return person(0, name = name, age = age)
}

type persons: []person
`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	env := ast.NewEnvironment()
	fmt.Println(prog.Transpile(env))
}
