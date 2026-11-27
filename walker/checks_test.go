package walker

import (
	"strings"
	"testing"

	"github.com/vapourlang/vapour/diagnostics"
	"github.com/vapourlang/vapour/lexer"
	"github.com/vapourlang/vapour/parser"
)

func walkTest(t *testing.T, code string) *Walker {
	t.Helper()
	l := lexer.NewTest(code)
	l.Run()
	p := parser.New(l)
	prog := p.Run()

	if p.HasError() {
		t.Fatalf("failed to parse: %v", p.Errors())
	}

	w := New()
	w.Run(prog)
	return w
}

func TestCustomInfix(t *testing.T) {
	w := walkTest(t, `let x: bool = 1 %in% (1, 2)
let y: int = (1, 2) %>% sum()

# should fail, %in% returns a logical
let z: char = 1 %in% (1, 2)
`)

	w.testDiagnostics(t, diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
	})
}

func TestCallbackReturn(t *testing.T) {
	w := walkTest(t, `type math: func(int, int) int

func apply_math(x: int, y: int, cb: math): int {
  stopifnot(!missing(x), !missing(y), !missing(cb))
  return cb(x, y)
}

apply_math(1, 2, (x: int, y: int): int => {
  return x * y
})

# should fail, wrong return type
apply_math(1, 2, (x: int, y: int): char => {
  stopifnot(!missing(x), !missing(y))
  return "a"
})
`)

	w.testDiagnostics(t, diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
	})
}

func TestRange(t *testing.T) {
	w := walkTest(t, `let n: int = 5
let a: int = n..10
let b: int = 1..n
let s: char = "a"

# should fail, char in range
let c: int = s..3

# should fail, char in range
let d: int = 1..s
`)

	w.testDiagnostics(t, diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
	})
}

func TestLoopVariable(t *testing.T) {
	w := walkTest(t, `let xs: char = ("a", "b")

for(let x: char in xs) {
  print(x)
}

for(let n: num in 1..10) {
  print(n)
}

# should fail, elements are char
for(let x: int in xs) {
  print(x)
}
`)

	w.testDiagnostics(t, diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
	})
}

func TestIterateListType(t *testing.T) {
	w := walkTest(t, `type user: object {
  name: char
}

type users: []user

let us: users = users(user(name = "a"))

for(let u: user in us) {
  print(u)
}

type nums: list { int }

let ns: nums = nums(1, 2)

for(let n: int in ns) {
  print(n)
}
`)

	w.testDiagnostics(t, diagnostics.Diagnostics{})
}

func TestLoopVariableListType(t *testing.T) {
	w := walkTest(t, `type user: object {
  name: char
}

type users: []user

let us: users = users(user(name = "a"))

# should fail, elements are user
for(let u: char in us) {
  print(u)
}
`)

	w.testDiagnostics(t, diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
	})

	if !strings.Contains(w.Errors()[0].Message, "elements are `user`") {
		t.Fatalf("expected element type error, got %v", w.Errors()[0].Message)
	}
}

func TestNestedAttributes(t *testing.T) {
	w := walkTest(t, `type address: object {
  city: char
}

type customer: object {
  name: char,
  home: address
}

let c: customer = customer(name = "a", home = address(city = "b"))

c$home$city = "Brussels"

let city: char = c$home$city

# should fail, city is a char
c$home$city = 1

# should fail, zip does not exist
c$home$zip = "1000"

# should fail, city is a char
let wrong: int = c$home$city
`)

	w.testDiagnostics(t, diagnostics.Diagnostics{
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
		{Severity: diagnostics.Fatal},
	})
}

func TestLoopControl(t *testing.T) {
	w := walkTest(t, `for(let i: int in 1..10) {
  if(i == 2) {
    next
  }
  if(i == 5) {
    break
  }
  print(i)
}

let count: int = 0
repeat {
  count += 1
  if(count >= 3) {
    break
  }
}
`)

	w.testDiagnostics(t, diagnostics.Diagnostics{})
}
