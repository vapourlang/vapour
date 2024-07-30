package parser

import (
	"fmt"
	"testing"

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

	fmt.Println(prog.String())
}

func TestFunc(t *testing.T) {
	code := `func add(x: int = 1, y: int = 2) int {
  let total: int = x + y * 2
  return total
} `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestPipe(t *testing.T) {
	code := `func add(): null {
  df |>
    mutate(x = 1)
} `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
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

	fmt.Println(prog.String())
}

func TestComment(t *testing.T) {
	code := `#' @return something
func add(): int | number {
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

	fmt.Println(prog.String())
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

	fmt.Println(prog.String())
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

	fmt.Println(prog.String())
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

	fmt.Println(prog.String())
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

	fmt.Println(prog.String())
}

func TestTElipsis(t *testing.T) {
	code := `func foo(...: any) char {
  paste0(..., collapse = ", ")
}  `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestS3(t *testing.T) {
	code := `
type person: struct {
  int | num,
  name: string,
  age: int
}

func (p person) getAge(): int {
  return p.age
}

func (p person) setAge(n: int): null {
  p.age = n
}

func create(name: string, age: int): person {
  return person(0, name = name, age = age)
}

type persons: []person

create(name = "hello")
`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestRange(t *testing.T) {
	code := `let x: int | na = 1..10
`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestFor(t *testing.T) {
	code := `for(let i:int = 1 in 1..nrow(df)) {
  print(i)
}

func foo(...: int): int {
  sum(...)
}

let x: int = (1, 20, 23) `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestWhile(t *testing.T) {
	code := `while(i < 10) {
  print(i)
}
`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestNamespace(t *testing.T) {
	code := `let x: dataframe = dplyr::mutate(cars, speed > 2) `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestIf(t *testing.T) {
	code := `let x: bool = (1,2,3)

if (x) {
  print("true")
} else {
  print("false")
}

func foo(n: int) null {
  # comment
  if(n == 1) {
    print(true)
  }
}
`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}
