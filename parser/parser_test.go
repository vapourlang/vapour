package parser

import (
	"fmt"
	"testing"

	"github.com/devOpifex/vapour/lexer"
)

func TestBasic(t *testing.T) {
	fmt.Println("---------------------------------------------------------- basic")
	code := `let x: int | num = 1  

let y: int

y = 2
`

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestFunc(t *testing.T) {
	fmt.Println("---------------------------------------------------------- func")
	code := `func add(x: int = 1, y: int = 2): int {
  let total: int = x + y * 2
  return total
} `

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestPipe(t *testing.T) {
	fmt.Println("---------------------------------------------------------- pipe")
	code := `func add(): null {
  df |>
    mutate(x = 1)
} `

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestString(t *testing.T) {
	fmt.Println("---------------------------------------------------------- string")
	code := `let x: string = "a \"string\""
let y: string = 'single quotes'`

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestComment(t *testing.T) {
	fmt.Println("---------------------------------------------------------- comment")
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
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestMethod(t *testing.T) {
	fmt.Println("---------------------------------------------------------- method")
	code := `func (o: obj) method(n: int): char {
  return "hello"
}`

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestTypeDeclaration(t *testing.T) {
	fmt.Println("---------------------------------------------------------- type declaration")
	code := `type userId: int

type obj: struct {
  int | string,
  name: string,
  id: int
} `

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestIdent(t *testing.T) {
	fmt.Println("---------------------------------------------------------- ident")
	code := `let x: int = (1,2,3)

x[1, 2] = 15

x[[3]] = 15

df.x = 23

print(x) `

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestTElipsis(t *testing.T) {
	fmt.Println("---------------------------------------------------------- ellipsis")
	code := `func foo(...: any): char {
  paste0(..., collapse = ", ")
}  `

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestS3(t *testing.T) {
	fmt.Println("---------------------------------------------------------- s3")
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
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestRange(t *testing.T) {
	fmt.Println("---------------------------------------------------------- range")
	code := `let x: int | na = 1..10
`

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestFor(t *testing.T) {
	fmt.Println("---------------------------------------------------------- for")
	code := `for(let i:int = 1 in 1..nrow(df)) {
  print(i)
}

func foo(...: int): int {
  sum(...)
}

let x: int = (1, 20, 23) `

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestWhile(t *testing.T) {
	fmt.Println("---------------------------------------------------------- while")
	code := `while(i < 10) {
  print(i)
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestNamespace(t *testing.T) {
	fmt.Println("---------------------------------------------------------- namespace")
	code := `let x: dataframe = cars |>
dplyr::mutate(speed > 2) `

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestIf(t *testing.T) {
	fmt.Println("---------------------------------------------------------- if")
	code := `let x: bool = (1,2,3)

if (x) {
  print(
	  "true",
		1
  )
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

	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestSquare(t *testing.T) {
	fmt.Println("---------------------------------------------------------- square")
	code := `let x: int = (1,2,3)

x[2] = 3

let y: int = list(1,2,3)

y[[1]] = 1

let zz: string = ("hello|world", "hello|again")
let z: char = strsplit(zz[2], "\\|")[[1]]
`

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestFunctionParam(t *testing.T) {
	fmt.Println("---------------------------------------------------------- function param")
	code := `func foo(fn: function = (x: int): int => {return x + 1}, y: int = 2): int {
  return sapply(x, fn) + y
}

func bar(
  x: int,
	y: int = 1
): int {
  return x + y
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	if p.HasError() {
		for _, e := range p.Errors() {
			fmt.Println(e)
		}
		return
	}

	fmt.Println(prog.String())
}

func TestCall(t *testing.T) {
	fmt.Println("---------------------------------------------------------- call")
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
	p := New(l)

	prog := p.Run()

	if p.HasError() {
		for _, e := range p.Errors() {
			fmt.Println(e)
		}
		return
	}

	fmt.Println(prog.String())
}

func TestNestedCall(t *testing.T) {
	fmt.Println("---------------------------------------------------------- nested call")
	code := `
x$val = list(
	list(
		arg = parts[1] |> trimws(),
		types = types |> trimws()
	)
)
`

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	if p.HasError() {
		for _, e := range p.Errors() {
			fmt.Println(e)
		}
		return
	}

	fmt.Println(prog.String())
}

func TestMissing(t *testing.T) {
	fmt.Println("---------------------------------------------------------- missing")
	code := `
# can be missing
func hello(what: char): char {
  if(missing(what)) {
	  what = "Vapour"
	}
  sprintf("hello, %s!", what)
}

hello()
`

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	p.Errors().Print()

	fmt.Println(prog.String())
}

func TestDecorators(t *testing.T) {
	fmt.Println("---------------------------------------------------------- decorator")
	code := `
@class(x, y, z)
type custom: list {
  x: char,
	id: int
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	if p.HasError() {
		for _, e := range p.Errors() {
			fmt.Println(e)
		}
		return
	}

	fmt.Println(prog.String())
}

func TestFnType(t *testing.T) {
	fmt.Println("---------------------------------------------------------- function type")
	code := `type state_fn: func(x: int | na, y: int): int 
`

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	if p.HasError() {
		for _, e := range p.Errors() {
			fmt.Println(e)
		}
		return
	}

	fmt.Println(prog.String())
}

func TestAnonymous(t *testing.T) {
	fmt.Println("---------------------------------------------------------- anonymous")
	code := `let y: int = (1,2,3)

const x: char = "world"
lapply(("hello", x), (z: char): null => {
  print(z)
})

lapply(1..10, (z: char): null => {
  print(z)
})

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
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}
