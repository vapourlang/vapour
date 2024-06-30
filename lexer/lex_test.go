package lexer

import (
	"fmt"
	"os"
	"testing"
)

var itemName = map[itemType]string{
	itemError:         "error",
	itemIdent:         "identifier",
	itemDoubleQuote:   "\"",
	itemSingleQuote:   "'",
	itemAssign:        "assign",
	itemLeftCurly:     "{",
	itemRightCurly:    "}",
	itemLeftParen:     "(",
	itemRightParen:    ")",
	itemString:        "string",
	itemInteger:       "integer",
	itemFloat:         "float",
	itemNamespace:     "::",
	itemMathOperation: "operation",
}

func print(l *lexer) {
	fmt.Fprintln(os.Stdout, "----")
	for _, v := range l.items {
		name := itemName[v.class]
		fmt.Fprintf(os.Stdout, "%v: %v %v\n", v.class, v.val, name)
	}
	fmt.Fprintln(os.Stdout, "----")
	fmt.Fprintf(os.Stdout, "lexed %v tokens\n", len(l.items))
}

func TestBasic(t *testing.T) {
	code := `x <- 1 + 2 - 1
					y <- 2 `

	l := &lexer{
		input: code,
	}

	l.run()

	if len(l.items) == 0 {
		t.Fatal("No items where lexed")
	}

	print(l)
}

func TestFunction(t *testing.T) {
	code := `foo <- function(x = 1) {
	x + 1
}

x <- foo(2)

print(x) `

	l := &lexer{
		input: code,
	}

	l.run()

	if len(l.items) == 0 {
		t.Fatal("No items where lexed")
	}

	print(l)
}

func TestObjects(t *testing.T) {
	code := `x <- list(1, 2)

y = c(2, 3)
z <- data.frame()

str <- "Hello, world!" `

	l := &lexer{
		input: code,
	}

	l.run()

	if len(l.items) == 0 {
		t.Fatal("No items where lexed")
	}

	print(l)
}
