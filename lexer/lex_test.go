package lexer

import (
	"testing"

	"github.com/vapourlang/vapour/token"
)

func TestDeclare(t *testing.T) {
	code := `let x: int | na = 1
const y: char = "hello"
`

	l := NewTest(code)

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	tokens :=
		[]token.ItemType{
			token.ItemLet,
			token.ItemIdent,
			token.ItemColon,
			token.ItemTypes,
			token.ItemOr,
			token.ItemTypes,
			token.ItemAssign,
			token.ItemInteger,
			token.ItemNewLine,
			token.ItemConst,
			token.ItemIdent,
			token.ItemColon,
			token.ItemTypes,
			token.ItemAssign,
			token.ItemDoubleQuote,
			token.ItemString,
			token.ItemDoubleQuote,
		}

	for i, token := range tokens {
		actual := l.Items[i].Class
		if actual != token {
			t.Fatalf(
				"token %v expected `%v`, got `%v`",
				i,
				token,
				actual,
			)
		}
	}
}

func TestSimpleTypes(t *testing.T) {
	code := `type userid: int
type something: char | null
`

	l := NewTest(code)

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	tokens :=
		[]token.ItemType{
			token.ItemTypesDecl,
			token.ItemTypes,
			token.ItemColon,
			token.ItemTypes,
			token.ItemNewLine,
			token.ItemTypesDecl,
			token.ItemTypes,
			token.ItemColon,
			token.ItemTypes,
			token.ItemOr,
			token.ItemTypes,
		}

	for i, token := range tokens {
		actual := l.Items[i].Class
		if actual != token {
			t.Fatalf(
				"token %v expected `%v`, got `%v`",
				i,
				token,
				actual,
			)
		}
	}
}

func TestObjectTypes(t *testing.T) {
	code := `type thing: object {
  id: int,
	name: char
}

type lst: list { num | na }

type df: dataframe {
  name: char,
	id: int
}

type multiple: []int
`

	l := NewTest(code)

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	tokens :=
		[]token.ItemType{
			token.ItemTypesDecl,
			token.ItemTypes,
			token.ItemColon,
			token.ItemObjObject,
			token.ItemLeftCurly,
			token.ItemNewLine,
			token.ItemIdent,
			token.ItemColon,
			token.ItemTypes,
			token.ItemComma,
			token.ItemNewLine,
			token.ItemIdent,
			token.ItemColon,
			token.ItemTypes,
			token.ItemNewLine,
			token.ItemRightCurly,
			token.ItemNewLine,
			token.ItemNewLine,
			token.ItemTypesDecl,
			token.ItemTypes,
			token.ItemColon,
		}

	for i, token := range tokens {
		actual := l.Items[i].Class
		if actual != token {
			t.Fatalf(
				"token %v expected `%v`, got `%v`",
				i,
				token,
				actual,
			)
		}
	}
}

func TestComment(t *testing.T) {
	code := `# this is a comment

# this is another comment
`

	l := NewTest(code)

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	tokens :=
		[]token.ItemType{
			token.ItemComment,
			token.ItemNewLine,
			token.ItemNewLine,
			token.ItemComment,
			token.ItemNewLine,
		}

	for i, token := range tokens {
		actual := l.Items[i].Class
		if actual != token {
			t.Fatalf(
				"token %v expected `%v`, got `%v`",
				i,
				token,
				actual,
			)
		}
	}
}

func TestCall(t *testing.T) {
	code := `print(1)

sum(1, 2.3, 3)

foo(x = 1, y = 2, 'hello')
`

	l := NewTest(code)

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	tokens :=
		[]token.ItemType{
			token.ItemIdent,
			token.ItemLeftParen,
			token.ItemInteger,
			token.ItemRightParen,
			token.ItemNewLine,
			token.ItemNewLine,
			token.ItemIdent,
			token.ItemLeftParen,
			token.ItemInteger,
			token.ItemComma,
			token.ItemFloat,
			token.ItemComma,
			token.ItemInteger,
			token.ItemRightParen,
			token.ItemNewLine,
			token.ItemNewLine,
			token.ItemIdent,
			token.ItemLeftParen,
			token.ItemIdent,
			token.ItemAssign,
			token.ItemInteger,
			token.ItemComma,
			token.ItemIdent,
			token.ItemAssign,
			token.ItemInteger,
			token.ItemComma,
			token.ItemSingleQuote,
			token.ItemString,
			token.ItemSingleQuote,
			token.ItemRightParen,
			token.ItemNewLine,
		}

	for i, token := range tokens {
		actual := l.Items[i].Class
		if actual != token {
			t.Fatalf(
				"token %v expected `%v`, got `%v`",
				i,
				token,
				actual,
			)
		}
	}
}

func TestForWhile(t *testing.T) {
	code := `for(let i: int in 1..10) {}

while(i < 10) {}
`

	l := NewTest(code)

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	tokens :=
		[]token.ItemType{
			token.ItemFor,
			token.ItemLeftParen,
			token.ItemLet,
			token.ItemIdent,
			token.ItemColon,
			token.ItemTypes,
			token.ItemIn,
			token.ItemInteger,
			token.ItemRange,
			token.ItemInteger,
			token.ItemRightParen,
			token.ItemLeftCurly,
			token.ItemRightCurly,
			token.ItemNewLine,
			token.ItemNewLine,
			token.ItemWhile,
			token.ItemLeftParen,
			token.ItemIdent,
			token.ItemLessThan,
			token.ItemInteger,
			token.ItemRightParen,
			token.ItemLeftCurly,
			token.ItemRightCurly,
		}

	for i, token := range tokens {
		actual := l.Items[i].Class
		if actual != token {
			t.Fatalf(
				"token %v expected `%v`, got `%v`",
				i,
				token,
				actual,
			)
		}
	}
}

func TestFunctionLiteral(t *testing.T) {
	code := `func foo(x: int, y: num = 1): num {
  return x + y
}
`

	l := NewTest(code)

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	tokens :=
		[]token.ItemType{
			token.ItemFunction,
			token.ItemIdent,
			token.ItemLeftParen,
			token.ItemIdent,
			token.ItemColon,
			token.ItemTypes,
			token.ItemComma,
			token.ItemIdent,
			token.ItemColon,
			token.ItemTypes,
			token.ItemAssign,
			token.ItemInteger,
			token.ItemRightParen,
			token.ItemColon,
			token.ItemTypes,
			token.ItemLeftCurly,
			token.ItemNewLine,
			token.ItemReturn,
			token.ItemIdent,
			token.ItemPlus,
			token.ItemIdent,
			token.ItemNewLine,
			token.ItemRightCurly,
			token.ItemNewLine,
		}

	for i, token := range tokens {
		actual := l.Items[i].Class
		if actual != token {
			t.Fatalf(
				"token %v expected `%v`, got `%v`",
				i,
				token,
				actual,
			)
		}
	}
}

func TestIf(t *testing.T) {
	code := `if(x > 2) {
  print(1)
} else if (TRUE) {
  # nothing
} else {
  # nothing
}
`

	l := NewTest(code)

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	tokens :=
		[]token.ItemType{
			token.ItemIf,
			token.ItemLeftParen,
			token.ItemIdent,
			token.ItemGreaterThan,
			token.ItemInteger,
			token.ItemRightParen,
			token.ItemLeftCurly,
			token.ItemNewLine,
			token.ItemIdent,
			token.ItemLeftParen,
			token.ItemInteger,
			token.ItemRightParen,
			token.ItemNewLine,
			token.ItemRightCurly,
			token.ItemElse,
			token.ItemIf,
			token.ItemLeftParen,
			token.ItemBool,
			token.ItemRightParen,
			token.ItemLeftCurly,
			token.ItemNewLine,
			token.ItemComment,
			token.ItemNewLine,
			token.ItemRightCurly,
			token.ItemElse,
			token.ItemLeftCurly,
			token.ItemNewLine,
			token.ItemComment,
			token.ItemNewLine,
			token.ItemRightCurly,
		}

	for i, token := range tokens {
		actual := l.Items[i].Class
		if actual != token {
			t.Fatalf(
				"token %v expected `%v`, got `%v`",
				i,
				token,
				actual,
			)
		}
	}
}

func TestIncrement(t *testing.T) {
	code := `let x: int = 1

x += 1
`

	l := NewTest(code)

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	tokens :=
		[]token.ItemType{
			token.ItemLet,
			token.ItemIdent,
			token.ItemColon,
			token.ItemTypes,
			token.ItemAssign,
			token.ItemInteger,
			token.ItemNewLine,
			token.ItemNewLine,
			token.ItemIdent,
			token.ItemAssignInc,
			token.ItemInteger,
		}

	for i, token := range tokens {
		actual := l.Items[i].Class
		if actual != token {
			t.Fatalf(
				"token %v expected `%v`, got `%v`",
				i,
				token,
				actual,
			)
		}
	}
}

func TestFuncType(t *testing.T) {
	code := `type fn: func(int, []num) int`

	l := NewTest(code)

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	tokens :=
		[]token.ItemType{
			token.ItemTypesDecl,
			token.ItemTypes,
			token.ItemColon,
			token.ItemObjFunc,
			token.ItemLeftParen,
			token.ItemTypes,
			token.ItemComma,
			token.ItemTypesList,
			token.ItemTypes,
			token.ItemRightParen,
			token.ItemTypes,
		}

	for i, token := range tokens {
		actual := l.Items[i].Class
		if actual != token {
			t.Fatalf(
				"token %v expected `%v`, got `%v`",
				i,
				token,
				actual,
			)
		}
	}
}

func TestTypeImport(t *testing.T) {
	code := `let x: tibble::tbl = as.tibble(cars)
let y: tibble::[]tbl | int`

	l := NewTest(code)

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	tokens :=
		[]token.ItemType{
			token.ItemLet,
			token.ItemIdent,
			token.ItemColon,
			token.ItemTypesPkg,
			token.ItemNamespace,
			token.ItemTypes,
			token.ItemAssign,
			token.ItemIdent,
			token.ItemLeftParen,
			token.ItemIdent,
			token.ItemRightParen,
			token.ItemNewLine,
			token.ItemLet,
			token.ItemIdent,
			token.ItemColon,
			token.ItemTypesPkg,
			token.ItemNamespace,
			token.ItemTypesList,
			token.ItemTypes,
			token.ItemOr,
			token.ItemTypes,
		}

	for i, token := range tokens {
		actual := l.Items[i].Class
		if actual != token {
			t.Fatalf(
				"token %v expected `%v`, got `%v`",
				i,
				token,
				actual,
			)
		}
	}
}
