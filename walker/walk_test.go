package walker

import (
	"fmt"
	"testing"

	"github.com/devOpifex/vapour/diagnostics"
	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
)

func (w *Walker) testDiagnostics(t *testing.T, expected diagnostics.Diagnostics) {
	if len(w.errors) != len(expected) {
		w.errors.Print()
		t.Fatalf(
			"expected %v diagnostics, got %v",
			len(expected),
			len(w.errors),
		)
	}

	for index, e := range expected {
		if e.Severity != w.errors[index].Severity {
			w.errors[index].Print()
			t.Fatalf(
				"diagnostics %v, expected severity %v, got %v",
				index,
				e.Severity,
				w.errors[index].Severity,
			)
		}
	}
}

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

	w.Run(prog)

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Info},
		{Severity: diagnostics.Info},
	}

	w.testDiagnostics(t, expected)
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

	w.Run(prog)

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
	}

	w.testDiagnostics(t, expected)
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
	p := parser.New(l)

	prog := p.Run()

	w := New()
	w.Run(prog)

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
	}

	w.testDiagnostics(t, expected)
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

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
	}

	w.testDiagnostics(t, expected)
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
	p := parser.New(l)

	prog := p.Run()

	w := New()
	w.Run(prog)

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
	}

	w.testDiagnostics(t, expected)
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

	w.Run(prog)

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Info},
	}

	w.testDiagnostics(t, expected)
}

func TestExists(t *testing.T) {
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
		l.Errors().Print()
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

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Hint},
		{Severity: diagnostics.Hint},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Info},
	}

	w.testDiagnostics(t, expected)
}

func TestTypeMatch(t *testing.T) {
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

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
	}

	w.testDiagnostics(t, expected)
}

func TestR(t *testing.T) {
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

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Hint},
		{Severity: diagnostics.Hint},
		{Severity: diagnostics.Hint},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
	}

	w.testDiagnostics(t, expected)
}

func TestFunction(t *testing.T) {
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

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
	}

	w.testDiagnostics(t, expected)
}

func TestSquare(t *testing.T) {
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

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
	}

	w.testDiagnostics(t, expected)
}

func TestBasic(t *testing.T) {
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

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Info},
		{Severity: diagnostics.Info},
	}

	w.testDiagnostics(t, expected)
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

	w.Run(prog)

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
	}

	w.testDiagnostics(t, expected)
}

func TestListTypes(t *testing.T) {
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

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Info},
	}

	w.testDiagnostics(t, expected)
}

func TestFor(t *testing.T) {
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

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
	}

	w.testDiagnostics(t, expected)
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

	w.Run(prog)

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Info},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
	}

	w.testDiagnostics(t, expected)
}

func TestIncrement(t *testing.T) {
	code := `
let x: int = 1

x += 1

# should fail
x += "aaah"
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	w.Run(prog)

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
	}

	w.testDiagnostics(t, expected)
}

func TestReal(t *testing.T) {
	code := `
type lst: list {
  int | num
}

# should fail, multiple types
type mat: matrix {
  int | num
}

@matrix(nrow = 2, ncol = 4, wrong = true)
type matty: matrix {
  int
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	w.Run(prog)

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Info},
		{Severity: diagnostics.Info},
		{Severity: diagnostics.Info},
	}

	w.testDiagnostics(t, expected)
}

func TestFactor(t *testing.T) {
	code := `
@factor(wrong = TRUE)
type fac: factor {
  int
}

type fct: factor {
  int | num
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	w.Run(prog)

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Info},
		{Severity: diagnostics.Info},
	}

	w.testDiagnostics(t, expected)
}

func TestSignature(t *testing.T) {
	code := `
type math: func(int, int) int

func apply_math(x: int = 2, y: int = 2, cb: math): int {
  return cb(x, y)
}

func multiply(x: int = 2, y: int = 1): int {
  return x * y
}

apply_math(1, 2, multiply)

apply_math(1, 2, (x: int, y: int): int => {
  return x + y
})

# should fail, wrong signatures
apply_math(1, 2, (x: int): int => {
  return x + 1
})

func zz(x: int = 2): num {
  return x + 2.1
}

apply_math(1, 2, zz)

apply_math(1, 2, (x: int, y: char): int => {
  return x + 1
})

func foo(x: int): math {
  return zz
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	w.Run(prog)

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Info},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Info},
	}

	w.testDiagnostics(t, expected)
}

func TestMethodCall(t *testing.T) {
	code := `
type rules: object {
  selector: char,
  rule: char
}

type linne: object {
  css: char,
  rules: []rules
}

#' @export
func create(): linne {
  return linne()
}

#' @export
func(l: linne) addRule(selector: char, ...: char): linne {
  l$rules <- append(l$rules, rule(selector = selector, rule = ""))
  return l
}

# error, already defined
func(l: linne) addRule(): linne {
  return l
}

# error, wrong type
addRule("wrongType", "hello")

# error, no args
addRule()

let l: linne = create()
addRule(l, "hello")
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	w.Run(prog)

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Info},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Info},
	}

	w.testDiagnostics(t, expected)
}
