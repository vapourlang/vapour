package lexer

import (
	"testing"
)

func TestFunction(t *testing.T) {
	code := `foo <- function(x = 1) {
	x + 1
}

x <- foo(2)

l.Print(x)`

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	l.Print()
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

	l.Print()
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

	l.Print()
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

	l.Print()
}

func TestComments(t *testing.T) {
	code := `# this is a function call
l.Print(cars)

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

	l.Print()
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

	l.Print()
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

	l.Print()
}

func TestIdentifier(t *testing.T) {
	code := `x <- 1
x <- data.frame(x = 1, y = 2)

my_function <- function(x) x + 1

my.function <- function(x) x - 1

l.Print(TRUE)`

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	l.Print()
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

	l.Print()
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

	l.Print()
}

func TestIf(t *testing.T) {
	code := `if(x | y) {
  l.Print("TRUE")
} else if (xx && yy) {
  break
} else {
  l.Print("FALSE")
}`

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	l.Print()
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

	l.Print()
}

func TestString(t *testing.T) {
	code := `x <- 'hello'

y <- 'world'

l.Print(paste(x, y))

long_str <- "hello, world!"

escaped <- "hello \"world\""`

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	l.Print()
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

	l.Print()
}

func TestMath(t *testing.T) {
	code := `x <- 1 + 1 / -2 * 3 ^ 2
if(2 %% 2){
  l.Print(TRUE)
}`

	l := &Lexer{
		Input: code,
	}

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	l.Print()
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

	l.Print()
}
