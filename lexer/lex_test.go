package lexer

import (
	"testing"
)

func TestBasicTypes(t *testing.T) {
	code := `type userId: int | null

const x: int = 1

let y: []int = list(1, 23, 33)

# structure(1..10, name = "", id = 0)
type item: struct {
  int
  # attributes
  category: string
}

item(1, category = "")

# structure(item,(), name = "", id = 0)
type nested: struct {
  item,
  # attributes
  name: string,
  id: int
}

nested(
  item(1..10, name = "hello", id = 1),
  category = "test"
)

# data.frame(name = ("a", "z"), id = 1..2)
type df: dataframe {
  name: string,
  id: int
}

df(name = "hello", id = 1)

# list(1, 2, 3)
type lst list {
  int | string
}

lst( 1,2 ,3)

# list(name = "hello", id = 1)
type obj: object {
  id: int,
  n: num
}

obj(
  id = 0,
  n = 3.14
)

# list(list(name = "hello", id = 1))
type objs []obj

objs(
  obj(),
  obj()
)

func foo(x: string = "hello"): string {
  return paste0(x, ", world")
}

func foo_bar(foo: fn = (x: string): string => paste0(x, 1))

let x: int = (1,3,4)

func (x obj) do(): string {
  paste0(x.v)
}

func foo(...: any): string {
  paste0(..., collapse = ", ")
}

for(let i: int = 1 in 1:10){
  print(i)
}

@class(x, y, z)
type custom: list {
  x: char,
	id: int
}

let df: any = data.frame(x = 1..3)
df$x$y = 2
`

	l := NewTest(code)

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	l.Print()
}
