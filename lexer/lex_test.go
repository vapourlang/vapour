package lexer

import (
	"fmt"
	"os"
	"testing"
)

var itemName = map[itemType]string{
	itemError:             "error",
	itemIdent:             "identifier",
	itemDoubleQuote:       "double quote",
	itemSingleQuote:       "single quote",
	itemAssign:            "assign",
	itemLeftCurly:         "curly left",
	itemRightCurly:        "curly right",
	itemLeftParen:         "paren left",
	itemRightParen:        "paren right",
	itemLeftSquare:        "square left",
	itemRightSquare:       "square right",
	itemString:            "string",
	itemInteger:           "integer",
	itemFloat:             "float",
	itemNamespace:         "namespace",
	itemMathOperation:     "operation",
	itemComment:           "comment",
	itemSpecialComment:    "special comment",
	itemRoxygenTagAt:      "@roxygen",
	itemRoxygenTag:        "roxygen tag",
	itemRoxygenTagContent: "roxygen content",
	itemTypeDef:           "type",
	itemTypeVar:           "type variable",
	itemDoubleEqual:       "double equal",
	itemLessThan:          "less than",
	itemGreaterThan:       "greater than",
	itemNotEqual:          "not equal",
	itemLessOrEqual:       "less or equal",
	itemGreaterOrEqual:    "greater or equal",
}

func print(l *lexer) {
	fmt.Fprintln(os.Stdout, "----")
	for _, v := range l.items {
		name := itemName[v.class]
		fmt.Fprintf(os.Stdout, "%v: %v [%v]\n", v.class, v.val, name)
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

str <- "Hello, world!"

x <- c("hello", "world") `

	l := &lexer{
		input: code,
	}

	l.run()

	if len(l.items) == 0 {
		t.Fatal("No items where lexed")
	}

	print(l)
}

func TestNamespace(t *testing.T) {
	code := `x <- dplyr::filter(cars, speed > 10L)`

	l := &lexer{
		input: code,
	}

	l.run()

	if len(l.items) == 0 {
		t.Fatal("No items where lexed")
	}

	print(l)
}

func TestClasses(t *testing.T) {
	code := `Person <- setRefClass("Person")
p <- Person$new()`

	l := &lexer{
		input: code,
	}

	l.run()

	if len(l.items) == 0 {
		t.Fatal("No items where lexed")
	}

	print(l)
}

func TestSquare(t *testing.T) {
	code := `x <- data.frame(x = 1:10, y = 1:10)
x[1, 1] <- 3L`

	l := &lexer{
		input: code,
	}

	l.run()

	if len(l.items) == 0 {
		t.Fatal("No items where lexed")
	}

	print(l)
}

func TestComments(t *testing.T) {
	code := `# this is a function call
print(cars)

x <- 1 # is equal to 1

#' This is a special comment
NULL

#' @param x An integer to add one to
#' @type x: numeric
#' @yield numeric
foo <- function(x) {
  x + 1
}

#' @param x Another integer to add one to
#' @type x: numeric | integer
#' @yield numeric | integer
foo <- function(x) {
  x + 1
}

# THIS SHOULD ERROR
#' @type x numeric | integer
foo <- function(x) {
  x + 1
}`

	l := &lexer{
		input: code,
	}

	l.run()

	if len(l.items) == 0 {
		t.Fatal("No items where lexed")
	}

	print(l)
}

func TestCompare(t *testing.T) {
	code := `x <- 1

x == 3
x != 2
x >= 1
x <= 2
x < 2
x > 2`

	l := &lexer{
		input: code,
	}

	l.run()

	if len(l.items) == 0 {
		t.Fatal("No items where lexed")
	}

	print(l)
}

func TestNumbers(t *testing.T) {
	code := `x <- 123
x <- 1.23
x <- 1.1
y <- 10^2
y <- 10^.2
y <- 10e2`

	l := &lexer{
		input: code,
	}

	l.run()

	if len(l.items) == 0 {
		t.Fatal("No items where lexed")
	}

	print(l)
}
