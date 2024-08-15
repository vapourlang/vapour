package transpiler

import (
	"fmt"
	"testing"

	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
)

func TestBasic(t *testing.T) {
	code := `let x: int | num = 1  `

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestFunc(t *testing.T) {
	code := `func add(x: int = 1, y: int = 2): int {
  let total: int = x + y * 2
  return total
} `

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestPipe(t *testing.T) {
	code := `func add(): null {
  df |>
    mutate(x = 1)
} `

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestString(t *testing.T) {
	code := `let x: string = "a \"string\""
let y: string = 'single quotes'`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
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

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestTypeDeclaration(t *testing.T) {
	code := `type userId: int

type obj: struct {
  int | string,
  name: string,
  id: int
} `

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestIdent(t *testing.T) {
	code := `let x: int = (1,2,3)

x[1, 2] = 15

x[[3]] = 15

df.x = 23

print(x) `

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestTElipsis(t *testing.T) {
	code := `func foo(...: any): char {
  paste0(..., collapse = ", ")
}  `

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestS3(t *testing.T) {
	code := `
type person: struct {
  int | num,
  name: string,
  age: int
}

func (p: person) getAge(): int {
  return p.age
}

func (p: person) setAge(n: int): null {
  p.age = n
}

func create(name: string, age: int): person {
  return person(0, name = name, age = age)
}

type persons: []person

create(name = "hello")
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestRange(t *testing.T) {
	code := `let x: int | na = 1..10
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestFor(t *testing.T) {
	code := `for(let i:int = 1 in 1..nrow(df)) {
  print(i)
}

func foo(...: int): int {
  sum(...)
}

let x: int = (1, 20, 23) `

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestWhile(t *testing.T) {
	code := `while(i < 10) {
  print(i)
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestNamespace(t *testing.T) {
	code := `let x: dataframe = cars |>
dplyr::mutate(speed > 2) `

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestIf(t *testing.T) {
	code := `let x: bool = (1,2,3)

if (x) {
  print("true")
} else {
  print("false")
}

func foo(n: int): null {
  # comment
  if(n == 1) {
    print(true)
  }
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestAnonymous(t *testing.T) {
	code := `let y: int = (1,2,3)

const x: char = "world"
lapply(("hello", x), (z: char): null => {
  print(z)
})

lapply(1..10, (z: char): null => {
  print(z)
})
 `

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestSquare(t *testing.T) {
	code := `let x: int = (1,2,3)

x[2] = 3

let y: int = list(1,2,3)

y[[1]] = 1

let zz: string = ("hello|world", "hello|again")
let z: char = strsplit(zz[2], "\\|")[[1]]
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestMethod(t *testing.T) {
	code := `func (o: obj) add(n: int): char {
  return "hello"
}

type person: struct{
  int,
	name: char
}

func (p: person) setName(name: char): null {
  p$name = 2
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestCall(t *testing.T) {
	code := `
bar(1, x = 2, "hello")

bar(
  1,
	x = 2,
	"hello"
)

foo(z = 2)

foo(1, 2, 3)

foo(
  z = "hello"
)

foo("hello")
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestDeclare(t *testing.T) {
	code := `let x: int

x = 2

type config: list {
  name: char,
	x: int
}

config(name = "hello")

# should fail, does not exist
let z: config = config(
  z = 2
)

z$name = 2
`

	fmt.Println("-----------------------------")
	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestType(t *testing.T) {
	code := `type person: struct {
  list,
	name: string
}

person(list(), name = "John")

# should fail, attr not in type
person(list(), age = 1)

person(1)

@class(x, y, z)
type cl: struct {
  int
}

let z: cl = cl(2)

@class(fr, lt)
type lst: list {
  int
}

let zzzz: lst = lst()
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestList(t *testing.T) {
	code := `
type person: list {
	name: char
}

type persons: []person

let peoples: persons = persons(
  person(name = "John"),
  person(name = "Jane")
)

type ints: []ints

let x: ints = ints(1,2,3)
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}
