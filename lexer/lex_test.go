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
	itemNamespaceInternal: "namespace internal",
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
	itemBool:              "boolean",
	itemDollar:            "dollar",
	itemComma:             "comma",
	itemColon:             "colon",
	itemQuestion:          "question mark",
	itemBacktick:          "backtick",
	itemInfix:             "infix",
	itemIf:                "if",
	itemBreak:             "break",
	itemElse:              "else",
	itemAnd:               "ampersand",
	itemOr:                "vertical bar",
	itemReturn:            "return",
	itemCCall:             "C/C++ call",
	itemNULL:              "null",
	itemNA:                "NA",
	itemNan:               "NaN",
	itemNACharacter:       "NA character",
	itemNAReal:            "NA real",
	itemNAComplex:         "NA complex",
	itemNAInteger:         "NA integer",
	itemPipe:              "native pipe",
	itemModulus:           "modulus",
	itemDoubleLeftSquare:  "double left square",
	itemDoubleRightSquare: "double right square",
	itemFor:               "for loop",
	itemRepeat:            "repeat",
	itemWhile:             "while loop",
	itemNext:              "next",
	itemIn:                "in",
	itemFunction:          "function",
}

func print(l *lexer) {
	fmt.Fprintln(os.Stdout, "\n----")
	for _, v := range l.items {
		name := itemName[v.class]
		fmt.Fprintf(os.Stdout, "%v: %v [%v]\n", v.class, v.val, name)
	}
	fmt.Fprintf(os.Stdout, "=====> lexed %v tokens\n", len(l.items))
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
	code := `x <- dplyr::filter(cars, speed > 10L)
pkg:::internal()`

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
	x <- 12.23
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

func TestIdentifier(t *testing.T) {
	code := `x <- 1
x <- data.frame(x = 1, y = 2)

my_function <- function(x) x + 1

my.function <- function(x) x - 1

print(TRUE)`

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
p <- Person$new()

x <- R6::R6Class(
  "Person",
  public = list(
    initialize = function(){}
  )
)

p2 <- x$new()

?dplyr::filter`

	l := &lexer{
		input: code,
	}

	l.run()

	if len(l.items) == 0 {
		t.Fatal("No items where lexed")
	}

	print(l)
}

func TestBacktick(t *testing.T) {
	code := "`%||%` <- function(lhs, rhs) lhs"

	l := &lexer{
		input: code,
	}

	l.run()

	if len(l.items) == 0 {
		t.Fatal("No items where lexed")
	}

	print(l)
}

func TestIf(t *testing.T) {
	code := `if(x | y) {
  print("TRUE")
} else if (xx && yy) {
  break
} else {
  print("FALSE")
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

func TestSpecialTypes(t *testing.T) {
	code := `x <- NULL

x <- Inf
x <- NA
x <- NaN
x <- NA_character_
x <- NA_complex_
x <- NA_real_
x <- NA_integer_`

	l := &lexer{
		input: code,
	}

	l.run()

	if len(l.items) == 0 {
		t.Fatal("No items where lexed")
	}

	print(l)
}

func TestString(t *testing.T) {
	code := `x <- 'hello'

y <- 'world'

print(paste(x, y))

long_str <- "hello, world!"

escaped <- "hello \"world\""`

	l := &lexer{
		input: code,
	}

	l.run()

	if len(l.items) == 0 {
		t.Fatal("No items where lexed")
	}

	print(l)
}

func TestPipe(t *testing.T) {
	code := `data |>
dplyr::mutate(x = x + 1)

data %>% filter(x < 2)

x %||% y`

	l := &lexer{
		input: code,
	}

	l.run()

	if len(l.items) == 0 {
		t.Fatal("No items where lexed")
	}

	print(l)
}

func TestMath(t *testing.T) {
	code := `x <- 1 + 1 + -2 * 3
if(2 %% 2){
  print(TRUE)
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

func TestLoop(t *testing.T) {
	code := `for(i in 1:10){
 if(i > 5)
    next
}

x <- 1

while(x < 10){
  x <- x + 1
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
