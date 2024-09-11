package transpiler

import (
	"testing"

	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
)

func (trans *Transpiler) testOutput(t *testing.T, expected string) {
	if trans.GetCode() == expected {
		return
	}
	t.Fatalf("expected:\n`%v`\ngot:\n`%v`", expected, trans.GetCode())
}

func TestBasic(t *testing.T) {
	code := `let x: int | num = 1  
const y: int = 1`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	expectations := `x = 1
y = 1

`
	trans.testOutput(t, expectations)
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

	expected := `add = function(x = 1,y = 2) {
total = x+y*2
return(total)
}`

	trans.testOutput(t, expected)
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

	expected := `add = function() {
df|>mutate(x=1
)
}`

	trans.testOutput(t, expected)
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

	expected := `x = "a \"char\""
y = 'single quotes'
`

	trans.testOutput(t, expected)
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

	expected := `#' @return something
add = function() {
# compute stuff
x = df|>mutate(x="hello"
, y=na
, b=TRUE
)|>select(x)
return(x)
}`

	trans.testOutput(t, expected)
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

	expected := `foo = function(...) {
paste0(..., collapse=", "
)
}`

	trans.testOutput(t, expected)
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

	expected := `x = 1:10
`

	trans.testOutput(t, expected)
}

func TestFor(t *testing.T) {
	code := `
for(let i: int in 1..nrow(df)) {
  print(i)
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	expected := `print(i)
`

	trans.testOutput(t, expected)
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

	expected := `while(i<10
) {
print(i)
}
`

	trans.testOutput(t, expected)
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

	expected := `x = cars|>dplyr::mutate(speed>2
)
`

	trans.testOutput(t, expected)
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

	expected := `x = c(1, 2, 3)
if(x){
print("true")
} else {
print("false")
}
foo = function(n) {
# comment
if(n==1
){
print(TRUE)
}}
`

	trans.testOutput(t, expected)
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

	expected := `y = c(1, 2, 3)
x = "world"
lapply(c("hello", x), function(z) {
print(z)
})
lapply(1:10
, function(z) {
print(z)
})
`

	trans.testOutput(t, expected)
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

	expected := `add.obj = function(o, n) {
return("hello")
}
setName.person = function(p, name) {
p$name=2
}
`

	trans.testOutput(t, expected)
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

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	expected := `x=2
structure(list(name="hello"
), class=c("config", "list"))
# should fail, does not exist
z = structure(list(z=2
), class=c("config", "list"))
z$name=2
`

	trans.testOutput(t, expected)
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

type math: func(int) int

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

	expected := `peoples = persons(structure(list(name="John"
), class=c("person", "list"))
, structure(list(name="Jane"
), class=c("person", "list"))
)
x = ints(1, 2, 3)
apply_math = function(vector,cb) {
return(cb(vector)
)
}
apply_math(c(1, 2, 3), function(x) {
return(x*3
)
})
`

	trans.testOutput(t, expected)
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

	expected := `x = c(1, 2, 3)
x[2]=3
y = list(1, 2, 3)
y[[1]]=1
zz = c("hello|world", "hello|again")
z = strsplit(zz[2
, ], "\\|")[[1]]
`

	trans.testOutput(t, expected)
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

	expected := `structure(list(1, "hello"), class=c("userid", "list"))
`

	trans.testOutput(t, expected)
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

	expected := `structure(42, name="xxx"
, class="st")
structure(list(name="hello"
), class=c("obj", "list"))
structure(list(name="hello"
), class=c("hello", "world")
structure(data.frame(name="hello"
, id=1
), names = c("name", "id"), class=c("df", "data.frame"))
`

	trans.testOutput(t, expected)
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

	expected := `foo = function(x) {
on.exit((function() {print("hello")
})())
return(1+1
)
}
`

	trans.testOutput(t, expected)
}

func TestStruct(t *testing.T) {
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

	expected := `create = function(name,age) {
return(structure(0, name=name
, age=age
, class="person")
)
}
create2 = function() {
return(structure(1, class="thing")
)
}
create3 = function() {
return(structure(2, class=c("more", "classes", "here")
)
}
`

	trans.testOutput(t, expected)
}

func TestObject(t *testing.T) {
	code := `
type df: dataframe {
  name: char,
	age: int
}

df(name = "hello", age = 1)

type thing: object {
  wheels: bool
}

thing(wheels = true)
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	expected := `structure(data.frame(name="hello"
, age=1
), names = c("name", "age"), class=c("df", "data.frame"))
structure(list(wheels=TRUE
), class=c("thing", "list"))
`

	trans.testOutput(t, expected)
}

func TestVector(t *testing.T) {
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

	expected := `c(3)
structure(list(1, "hello"), class=c("lst", "list"))
`

	trans.testOutput(t, expected)
}

func TestType(t *testing.T) {
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

	expected := `structure(list()
, name="John"
, class="person")
# should fail, attr not in type
structure(list()
, age=1
, class="person")
structure(1, class="person")
z = structure(2, class=c("x", "y", "z")
zzzz = structure(list(), class=c("fr", "lt")
set_age = function(p, age) {UseMethod("set_age")}
set_age.default = function(p, age) {
stop("not implemented")
}
`

	trans.testOutput(t, expected)
}

func TestIncrement(t *testing.T) {
	code := `
let x: int = 10

x += 2
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	expected := `x = 10
x=x+2
`

	trans.testOutput(t, expected)
}

func TestMatrix(t *testing.T) {
	code := `
@matrix(nrow = 2, ncol = 4)
type mat: matrix {
  int
}

mat((1, 2, 3))
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	expected := `structure(matrix(c(1, 2, 3), nrow=2
ncol=4
), class=c("mat", "matrix"))
`

	trans.testOutput(t, expected)
}

func TestFactor(t *testing.T) {
	code := `
@factor(levels = TRUE)
type fac: factor {
  int
}

fac((1, 2, 3))
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	expected := `structure(factor(c(1, 2, 3), levels=TRUE
), class=c("fac", "factor"))
`

	trans.testOutput(t, expected)
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

	expected := `bar(1, x=2
, "hello")
bar(1, x=2
, "hello")
foo(z=2
)
foo(1, 2, 3)
foo(z="hello"
)
foo("hello")
`

	trans.testOutput(t, expected)
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

	expected := `x = c(1, 2, 3)
x[1,2]=15
x[[3]]=15
df$x=23
print(x)
`

	trans.testOutput(t, expected)
}

func TestAttribute(t *testing.T) {
	code := `type person: object {
  age: int  
}

p$age = 2
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	expected := ``

	trans.testOutput(t, expected)
}
