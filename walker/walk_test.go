package walker

import (
	"fmt"
	"testing"

	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
)

func TestEnvironment(t *testing.T) {
	code := `
let z: int = 1

func addz(n: int, y: int): int | na {
	if(n == 1){
		return NA
	}

	return n + y
}

# should fail, this can be na
let result: int = addz(1, 2)

# should fail, comparing wrong types
if (1 == "hello") {
  print("1")
}

if (1 > 2.1) {
  print("1")
}

const y: int = 1

# should fail, is constant
y = 2
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("----------------------------- Env")
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestInfix(t *testing.T) {
	code := `let x: char = "hello"

# should fail, cannot be NA
x = NA

# should fail, types do not match
let z: char = 1

func add(n: int, y: int): int | na {
	if(n == 1){
		return NA
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

	fmt.Println("----------------------------- infix")
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestNamespace(t *testing.T) {
	code := `
# should fail, duplicated params
func foo(x: int, x: int): int {return x + y}

# should fail, duplicated params
func (x: int) bar(x: int): int {return x + y}
`

	l := lexer.NewTest(code)

	l.Run()
	fmt.Println("----------------------------- ns")
	p := parser.New(l)

	prog := p.Run()
	if p.HasError() {
		for _, e := range p.Errors() {
			fmt.Println(e)
		}
		return
	}

	w := New()
	w.Run(prog)

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

let s: int = sum(1,2,3)
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	w.Run(prog)

	fmt.Println("----------------------------- number")
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
lg("hello", something = 1)
`

	l := lexer.NewTest(code)

	l.Run()
	fmt.Println("----------------------------- call")
	p := parser.New(l)

	prog := p.Run()

	w := New()
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestMissing(t *testing.T) {
	code := `
# should warn, can be missing
func hello(what: char): char {
  sprintf("hello, %s!", what)
}

type dataset: dataframe {
  name: char
}

# should warn, can be missing
func h(dat: dataset): char {
  dat$name = "hello"
	return "done"
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("----------------------------- missing")
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestExists(t *testing.T) {
	fmt.Println("----------------------------- exists")
	code := `
# should fail, x does not exist
x = 1

pkg::fn(x = 2)

dplyr::filter(x = 2)

# should fail, z does not exist
func foo(y: int): int {
 return z
}
`

	l := lexer.NewTest(code)

	l.Run()
	if l.HasError() {
		fmt.Printf("lexer errored")
		l.Errors.Print()
		return
	}
	p := parser.New(l)

	prog := p.Run()
	if p.HasError() {
		fmt.Printf("parser errored")
		p.Errors().Print()
	}

	w := New()

	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestTypeMatch(t *testing.T) {
	fmt.Println("----------------------------- typematch")
	code := `
type userid: int

let me: userid = userid(1)

# should fail, wrong type
let him: userid = "hello"

# should fail, named
me = userid(x = 1)

type lst: list { int | na }

let theList: lst = lst(1, 2)

# should fail, wrong type
theList = lst("aaaa", 1)

type config: struct {
  char,
	x: int
}

type inline: object {
  first: int,
	second: char
}

config(2, x = 2)

# should fail, must be named
inline(1)

# should fail, first arg of struct cannot be named
config(u = 2, x = 2, z = 2)

# should fail, struct attribute must be named
config(2, 2)

# should fail, does not exist
inline(
  z = 2
)
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	if p.HasError() {
		fmt.Println(p.Errors())
		return
	}

	w := New()

	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestR(t *testing.T) {
	fmt.Println("----------------------------- R")
	code := `
# should fail, package not installed
xxx::foo()

dplyr::wrong_function()

let x: int = 1

if x == 1 {
  x <- 2
}

# should fail does not exist
y <- 2
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestFunction(t *testing.T) {
	fmt.Println("----------------------------- function")
	code := `
func foo(n: int): int {
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

  # should fail, returns does not exist
  return u
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestSquare(t *testing.T) {
	fmt.Println("----------------------------- square")
	code := `let x: int = (1,2,3)

x[2] = 3

let y: int = list(1,2,3)

y[[1]] = 1

let zz: char = ("hello|world", "hello|again")
let z: char = strsplit(zz[2], "\\|")[[1]]

x[1, 2] = 15

x[[3]] = 15

type xx: dataframe {
  name: int
}

let df: xx = xx(name = 1)

df$name = "hello"

# should fail, not generic
func (p: any) meth(): null {}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestBasic(t *testing.T) {
	fmt.Println("----------------------------- Basic")
	code := `let x: int | na = 1

x = 2

# should fail, it's already declared
let x: char = "hello"

type ids: struct {
  int,
	name: char
}

# should fail, already defined
type id: int

# should fail, cannot find type
let z: undefinedType = "hello"

# should fail, different types
let v: int = (10, "hello", NA)

# should fail, type mismatch
let wrongType: num = "hello"

if(xx == 1) {
	let x: int = 2
}

# should fail wrong types
let x: int = "hello" + 2

# should fail, x is int, expression coerces to num
x = 1 + 1.2

let Z: num = 1.2 + 3

# should fail, does not exist
x = 2

# should fail, does not exist
uu = "char"
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	if p.HasError() {
		p.Errors().Print()
		return
	}

	w := New()

	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestRecursiveTypes(t *testing.T) {
	code := `
type userid: int

type user: struct {
  userid,
	name: char
}

func create(id: userid): user {
  return user(id)
}

create(2)

type person: struct {
  char,
	name: string
}

# should fail, wrong type
person(2)
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("----------------------------- Recurse types")
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestListTypes(t *testing.T) {
	fmt.Println("----------------------------- list types")
	code := `
type userid: int

type user: object {
	name: char
}

type users: []user

#should fail, wrong type
let z: users = users(
  user(name = "john"),
	4
)

# should fail, named
let w: users = users(
  user(name = "john"),
  x = user(name = "john"),
)
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestFor(t *testing.T) {
	fmt.Println("----------------------------- for")
	code := `
type userid: int

let x: userid = 1

for(let i: int in x..10) {
  print(i)
}

type x: struct {
  int
}

let y: x = x(2)

# should fail, range has char..int
for(let i: int in y) {
  print(i)
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestUnused(t *testing.T) {
	code := `
# should warn, x is never used
func foo(x: int): int {
  return 1
}

# should warn, does not exist
print(y)

# should warn, might be missing
func bar(x: int): int {
  return x
}

func baz(x: int): int { 
  stopifnot(!missing(x))
  return x
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	fmt.Println("----------------------------- unused")
	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}

func TestReal(t *testing.T) {
	fmt.Println("----------------------------- real")
	code := `type man: object {
  name: char
}

func foo(x: int = 1): man | null {
  if(x > 10){
	  return NULL
	}

	return 2
}

let x: int = 1
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	w.Run(prog)

	if len(w.errors) > 0 {
		w.errors.Print()
		return
	}
}
