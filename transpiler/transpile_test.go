package transpiler

import (
	"fmt"
	"testing"

	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
)

func TestBasic(t *testing.T) {
	code := `let x: int | num = 1  
const y: int = 1`

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
	code := `let x: char = "a \"char\""
let y: char = 'single quotes'`

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

func TestIdent(t *testing.T) {
	code := `let x: int = (1,2,3)

x[1, 2] = 15

x[[3]] = 15

df$x = 23

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

func TestRange(t *testing.T) {
	fmt.Println("-------------------------------------------- range")
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
	fmt.Println("-------------------------------------------- for")
	code := `for(let i:int in 1..nrow(df)) {
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

type config: object {
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

type math: func(x: int): int

func apply_math(vector: int, cb: math): int {
  return cb(vector)
}

apply_math((1, 2, 3), (x: int): int => {
  return x * 3
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

let zz: char = ("hello|world", "hello|again")
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

func TestClass(t *testing.T) {
	code := `
type userid: object {
  id: int,
	name: char
}

userid(1, "hello")
`

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

type st: struct {
  int | char,
  name: char,
  id: int
}

st(42, name = "xxx")

type obj: object {
  name: char,
  id: int
}

obj(name = "hello")

@class(hello, world)
type thing: object {
  name: char
}

thing(name = "hello")

type df: dataframe {
  name: char,
	id: int
}

df(name = "hello", id = 1)
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestDefer(t *testing.T) {
	code := `
func foo(x: int): int {
	defer (): null => {print("hello")}
  return 1 + 1
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

func TestStruct(t *testing.T) {
	fmt.Println("-------------------------------------------- struct")
	code := `
type person: struct {
  int | num,
  name: char,
  age: int
}

func create(name: char, age: int): person {
  return person(0, name = name, age = age)
}

type thing: struct {
  int
}

func create2(): thing {
  return thing(1)
}

@class(more, classes, here)
type stuff: struct {
  int
}

func create3(): thing {
  return stuff(2)
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

func TestObject(t *testing.T) {
	fmt.Println("-------------------------------------------- object dataframe")
	code := `
type df: dataframe {
  name: char,
	age: int
}

df(name = "hello", age = 1)

type thing: object {
  wheels: bool
}

thing(wheels = TRUE)
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestVector(t *testing.T) {
	fmt.Println("-------------------------------------------- vector and list")
	code := `
type userid: int

userid(3)

type lst: list {
  int | char | na
}

lst(1, "hello")
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestType(t *testing.T) {
	fmt.Println("-------------------------------------------- type")
	code := `
type person: struct {
  list,
	name: char
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

@generic
func (p: any) set_age(age: int): any

@default
func (p: any) set_age(age: int): null {
  stop("not implemented")
}

type person: struct {
  char
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
