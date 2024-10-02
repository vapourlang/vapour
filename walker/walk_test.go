package walker

import (
	"fmt"
	"testing"

	"github.com/vapourlang/vapour/diagnostics"
	"github.com/vapourlang/vapour/environment"
	"github.com/vapourlang/vapour/lexer"
	"github.com/vapourlang/vapour/parser"
	"github.com/vapourlang/vapour/r"
)

func (w *Walker) testDiagnostics(t *testing.T, expected diagnostics.Diagnostics) {
	if len(w.Errors()) != len(expected) {
		w.Errors().Print()
		t.Fatalf(
			"expected %v diagnostics, got %v",
			len(expected),
			len(w.Errors()),
		)
	}

	for index, e := range expected {
		if e.Severity != w.Errors()[index].Severity {
			fmt.Printf("Error at %v\n", index)
			fmt.Printf("%v", w.Errors()[index])
			t.Fatalf(
				"diagnostics %v, expected severity %v, got %v",
				index,
				e.Severity,
				w.Errors()[index].Severity,
			)
		}
	}
}

func TestInfix(t *testing.T) {
	code := `let x: char = "hello"

# no longer fails on v0.0.5
x = NA

# should fail, types do not match
let z: char = 1

# should be fine since v0.0.5
func add(n: int = 1, y: int = 2): int {
	if(n == 1){
		return NA
	}

  return n + y
}

# no longer fails since v0.0.5
let result: int = add(1, 2)

# should fail, const must have single type
const v: int | char = 1

const c: int = 1

# should fail, it's a constant
c += 2
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
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
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
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
	}

	w.testDiagnostics(t, expected)
}

func TestExists(t *testing.T) {
	code := `
# should fail, x does not exist
x = 1

# package not installed
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

# should fail, vector expects unnamed args
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
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
	}

	w.testDiagnostics(t, expected)
}

func TestR(t *testing.T) {
	code := `
# should fail, package not installed
xxx::foo()

# should fail, wrong function
dplyr::wrong_function()

let x: int = 1

if(x == 1){
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

# should fail, missing return
lapply(1..10, (x: int): int => {})

# should fail, string of identifiers
lapply(1..10, (): int => {
	this is not valid but doesnt throw an error
})
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
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Info},
		{Severity: diagnostics.Fatal},
	}

	w.testDiagnostics(t, expected)
}

func TestSquare(t *testing.T) {
	code := `let x: int = (1,2,3)

# we no longer check these
x[2] = 3

let zz: char = ("hello|world", "hello|again")
let z: any = strsplit(zz[2], "\\|")[[1]]

x[1, 2] = 15

x[[3]] = 15

type xx: dataframe {
  name: int
}

let df: xx = xx(name = 1)

# shoudl fail
x$sth

# should fail, wrong type
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
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
	}

	w.testDiagnostics(t, expected)
}

func TestBasic(t *testing.T) {
	code := `let x: int = 1

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

# should warn, might be missing
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
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Info},
	}

	w.testDiagnostics(t, expected)
}

func TestFor(t *testing.T) {
	code := `
type userid: int

let x: userid = 10

for(let i: int in 1..x) {
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

# should fail, wrong arg name
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
	}

	w.testDiagnostics(t, expected)
}

func TestFactor(t *testing.T) {
	code := `
# should error, wrong arg
@factor(wrong = TRUE)
type fac: factor {
  int
}

# should error, multiple types
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

# should fail, method on any
func(l: any) addRule(): null {}
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
		{Severity: diagnostics.Info},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Info},
	}

	w.testDiagnostics(t, expected)
}

func TestTypeInfix(t *testing.T) {
	code := `
let x: int = 10

# cannot use $ on int
# wrong unknow on x
x$wrong = 2

type person: object {
  name: char
}

let p: person = person()

# should fail, wrong type
p$name = 2
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
		{Severity: diagnostics.Fatal},
	}

	w.testDiagnostics(t, expected)
}

func TestType(t *testing.T) {
	code := `
# should fail, duplicated attribute
type person: object {
  name: char,
	name: int
}

# should fail, multiple types
type fct: factor { num | int }
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
	}

	w.testDiagnostics(t, expected)
}

func TestDate(t *testing.T) {
	code := `
type dt: date = as.Date("2022-02-01")

let x: dt = date("2022-02-01")

type person: object {
  dob: date
}

# should fail wrong type
let p: person = person(dob = 2)

# should fail wrong type
p$dob = 1
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
	}

	w.testDiagnostics(t, expected)
}

func TestReturn(t *testing.T) {
	code := `
# should fail, missing return
func foo(x: int): int {}

func (x: int) foo(y: char): null {}

# should fail, wrong return type
func bar(x: int): int {
  return "hello"
}

func baz(): any {}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	w.Run(prog)

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Info},
		{Severity: diagnostics.Info},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Info},
	}

	w.testDiagnostics(t, expected)
}

func TestDecorators(t *testing.T) {
	code := `
@generic
func (p: person) foo(x: int): any

# fail already defined
@generic
func (p: any) foo(x: int): any

@generic
func (p: any) bar(x: int): person
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

func TestLibrary(t *testing.T) {
	code := `
library(something)

require(package)
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
	}

	w.testDiagnostics(t, expected)
}

func TestAccessor(t *testing.T) {
	code := `
let x: any = 1

# should WORK
x["hello"]
x[["hello"]]
x$hello
print(x$hello)

let globals: any = new.env(env = parent.env(), hash = TRUE)

let y: int = (1 ,2 ,3)
y = y[3]
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	w.Run(prog)

	expected := diagnostics.Diagnostics{}

	w.testDiagnostics(t, expected)
}

func TestEnvironment(t *testing.T) {
	code := `
let z: int = 1

func addz(n: int = 1, y: int = 2): int {
	if(n == 1){
		return NA
	}

	return n + y + z
}

# should fail, comparing wrong types
if (1 == "hello") {
  print("1")
}

if (1 > 2.1) {
  print("1")
}

const c: int = 1

# should fail, is constant
c = 2

print(c)
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	w := New()

	w.Run(prog)

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Info},
		{Severity: diagnostics.Fatal},
	}

	w.testDiagnostics(t, expected)
}

func TestTypeImport(t *testing.T) {
	code := `
let x: vape::user = vape::user(id = 1)

let y: vape::user = vape::user(id = "char")
let w: vape::user = vape::user(wrong = 2)

type custom: object {
  user: vape::user
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	environment.SetLibrary(r.LibPath())
	w := New()

	w.Run(prog)

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Info},
	}

	w.testDiagnostics(t, expected)
}

func TestTypeNA(t *testing.T) {
	code := `
let x: int = 2

x = NA
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	environment.SetLibrary(r.LibPath())
	w := New()

	w.Run(prog)

	expected := diagnostics.Diagnostics{}

	w.testDiagnostics(t, expected)
}

func TestTypeEnvironment(t *testing.T) {
	code := `
type e: environment {
  name: char,
	x: int
}

e(name = "hello")

# should fail, unknown attribute
e(z = 2)

# should fail, wrong type
e(x = true)
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	environment.SetLibrary(r.LibPath())
	w := New()

	w.Run(prog)

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
	}

	w.testDiagnostics(t, expected)
}

func TestSymbols(t *testing.T) {
	code := `
# should fail, type does not exist
let foo: bar = baz

# should fail, baz not found
foo = baz
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	environment.SetLibrary(r.LibPath())
	w := New()

	w.Run(prog)

	expected := diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Warn},
	}

	w.testDiagnostics(t, expected)
}
