package lexer

import (
	"fmt"
	"os"
	"testing"
)

var ItemName = map[ItemType]string{
	ItemError:             "error",
	ItemIdent:             "identifier",
	ItemDoubleQuote:       "double quote",
	ItemSingleQuote:       "single quote",
	ItemAssign:            "assign",
	ItemLeftCurly:         "curly left",
	ItemRightCurly:        "curly right",
	ItemLeftParen:         "paren left",
	ItemRightParen:        "paren right",
	ItemLeftSquare:        "square left",
	ItemRightSquare:       "square right",
	ItemString:            "string",
	ItemInteger:           "integer",
	ItemFloat:             "float",
	ItemNamespace:         "namespace",
	ItemNamespaceInternal: "namespace internal",
	ItemComment:           "comment",
	ItemSpecialComment:    "special comment",
	ItemRoxygenTagAt:      "@roxygen",
	ItemRoxygenTag:        "roxygen tag",
	ItemRoxygenTagContent: "roxygen content",
	ItemTypeDef:           "type",
	ItemTypeVar:           "type variable",
	ItemDoubleEqual:       "double equal",
	ItemLessThan:          "less than",
	ItemGreaterThan:       "greater than",
	ItemNotEqual:          "not equal",
	ItemLessOrEqual:       "less or equal",
	ItemGreaterOrEqual:    "greater or equal",
	ItemBool:              "boolean",
	ItemDollar:            "dollar",
	ItemComma:             "comma",
	ItemColon:             "colon",
	ItemSemiColon:         "semicolon",
	ItemQuestion:          "question mark",
	ItemBacktick:          "backtick",
	ItemInfix:             "infix",
	ItemIf:                "if",
	ItemBreak:             "break",
	ItemElse:              "else",
	ItemAnd:               "ampersand",
	ItemOr:                "vertical bar",
	ItemReturn:            "return",
	ItemC:                 "C call",
	ItemCall:              "C++ call",
	ItemFortran:           "Fortan call",
	ItemNULL:              "null",
	ItemNA:                "NA",
	ItemNan:               "NaN",
	ItemNACharacter:       "NA character",
	ItemNAReal:            "NA real",
	ItemNAComplex:         "NA complex",
	ItemNAInteger:         "NA integer",
	ItemPipe:              "native pipe",
	ItemModulus:           "modulus",
	ItemDoubleLeftSquare:  "double left square",
	ItemDoubleRightSquare: "double right square",
	ItemFor:               "for loop",
	ItemRepeat:            "repeat",
	ItemWhile:             "while loop",
	ItemNext:              "next",
	ItemIn:                "in",
	ItemFunction:          "function",
	ItemPlus:              "plus",
	ItemMinus:             "minus",
	ItemMultiply:          "multiply",
	ItemDivide:            "divide",
	ItemPower:             "power",
	ItemEOL:               "end of line",
}

func print(l *Lexer) {
	fmt.Fprintln(os.Stdout, "\n----")
	for _, v := range l.Items {
		name := ItemName[v.class]
		fmt.Fprintf(os.Stdout, "%v: %v [%v]\n", v.class, v.val, name)
	}
	fmt.Fprintf(os.Stdout, "=====> lexed %v tokens\n", len(l.Items))
}

func TestFunction(t *testing.T) {
	code := `foo <- function(x = 1) {
	x + 1
}

x <- foo(2)

print(x) `

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	print(l)
}

func TestObjects(t *testing.T) {
	code := `x <- list(1, 2)

y = c(2, 3)
z <- data.frame()

str <- "Hello, world!"

x <- c("hello", "world") `

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	print(l)
}

func TestNamespace(t *testing.T) {
	code := `x <- dplyr::filter(cars, speed > 10L)
pkg:::internal()`

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	print(l)
}

func TestSquare(t *testing.T) {
	code := `x <- data.frame(x = 1:10, y = 1:10)
x[1, 1] <- 3L`

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
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

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
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

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
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

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	print(l)
}

func TestIdentifier(t *testing.T) {
	code := `x <- 1
x <- data.frame(x = 1, y = 2)

my_function <- function(x) x + 1

my.function <- function(x) x - 1

print(TRUE)`

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
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

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	print(l)
}

func TestBacktick(t *testing.T) {
	code := "`%||%` <- function(lhs, rhs) lhs"

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
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

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
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

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	print(l)
}

func TestString(t *testing.T) {
	code := `x <- 'hello'

y <- 'world'

print(paste(x, y))

long_str <- "hello, world!"

escaped <- "hello \"world\""`

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	print(l)
}

func TestPipe(t *testing.T) {
	code := `data |>
dplyr::mutate(x = x + 1)

data %>% filter(x < 2)

x %||% y`

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	print(l)
}

func TestMath(t *testing.T) {
	code := `x <- 1 + 1 / -2 * 3 ^ 2
if(2 %% 2){
  print(TRUE)
}`

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
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

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	print(l)
}
