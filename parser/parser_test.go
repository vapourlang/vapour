package parser

import (
	"fmt"
	"testing"

	"github.com/devOpifex/vapour/lexer"
)

func TestBasic(t *testing.T) {
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
	code := `let x: string = "a \"string\""
let y: string = 'single quotes'`

	l := lexer.NewTest(code)

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

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestMethod(t *testing.T) {
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
	code := `let x: int | na = 1..10
`

	l := lexer.NewTest(code)

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

	l := lexer.NewTest(code)

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

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestNamespace(t *testing.T) {
	code := `let x: dataframe = cars |>
dplyr::mutate(speed > 2) `

	l := lexer.NewTest(code)

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
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
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
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestFunctionParam(t *testing.T) {
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
	code := `
bar(1, x = 2, "hello")

bar(
  1,
	x = 2,
	"hello"
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
