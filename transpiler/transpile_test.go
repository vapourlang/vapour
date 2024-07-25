package transpiler

import (
	"fmt"
	"testing"

	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
)

func TestBasic(t *testing.T) {
	code := `let x: int | num = 1  `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()

	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestFunc(t *testing.T) {
	fmt.Println("+++++++++++++++")
	code := `func add(x: int = 1, y: int = 2) int {
  let total: int = x + y * 2
  return total
}`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()

	trans.Transpile(prog)

	if len(trans.errors) > 0 {
		trans.errors.Print()
		return
	}

	fmt.Println(trans.GetCode())
}

func TestFail(t *testing.T) {
	fmt.Println("+++++++++++++++")
	code := `let x: int = 1

x = 2

# should fail, it's already declared
let x: string = "hello"

const y: int = 1

# should fail, it's a const
y = 2

type id: struct {
  int,
  name: string
}

type id: int
let id: number = 1

id(1, name = "hello")
`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()

	trans.Transpile(prog)

	if len(trans.errors) > 0 {
		trans.errors.Print()
		return
	}

	fmt.Println(trans.GetCode())
}

func TestType(t *testing.T) {
	fmt.Println("+++++++++++++++")
	code := `type obj: struct {
  int,
  name: string
}

let x: obj = obj(1, name = "hello")

type df: dataframe {
  name: string,
  id: int
}

let y: df = df(name = ("hello", "world"), id = (1,2))

type id: int
`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()

	trans.Transpile(prog)

	if len(trans.errors) > 0 {
		trans.errors.Print()
		return
	}

	fmt.Println(trans.GetCode())
}