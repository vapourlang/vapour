package parser

import (
	"fmt"
	"testing"

	"github.com/devOpifex/vapour/lexer"
)

func TestBasic(t *testing.T) {
	code := `let x: int | num = 1  `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestFunc(t *testing.T) {
	code := `func add(x: int = 1, y: int = 2): int {
  let total: int = x + y * 2
  return total
} `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestPipe(t *testing.T) {
	code := `func add(): null {
  df |>
    mutate(x = 1)
} `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestString(t *testing.T) {
	code := `let x: string = "a \"string\""
let y: string = 'single quotes'`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
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

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestMethod(t *testing.T) {
	code := `func (o obj) add(n: int): char {
  return "hello"
}`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestTypeDeclaration(t *testing.T) {
	code := `type userId: int

type obj: struct {
  int | string,
  name: string,
  id: int
} `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestIdent(t *testing.T) {
	code := `let x: int = (1,2,3)

x[1, 2] = 15

x[[3]] = 15

df.x = 23

print(x) `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestTElipsis(t *testing.T) {
	code := `func foo(...: any) char {
  paste0(..., collapse = ", ")
}  `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestS3(t *testing.T) {
	code := `
type person: struct {
  int | num,
  name: string,
  age: int
}

func (p person) getAge(): int {
  return p.age
}

func (p person) setAge(n: int): null {
  p.age = n
}

func create(name: string, age: int): person {
  return person(0, name = name, age = age)
}

type persons: []person

create(name = "hello")
`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestRange(t *testing.T) {
	code := `let x: int | na = 1..10
`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestFor(t *testing.T) {
	code := `for(let i:int = 1 in 1..nrow(df)) {
  print(i)
}

func foo(...: int): int {
  sum(...)
}

let x: int = (1, 20, 23) `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestWhile(t *testing.T) {
	code := `while(i < 10) {
  print(i)
}
`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestNamespace(t *testing.T) {
	code := `let x: dataframe = cars |>
dplyr::mutate(speed > 2) `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestIf(t *testing.T) {
	code := `let x: bool = (1,2,3)

if (x) {
  print("true")
} else {
  print("false")
}

func foo(n: int) null {
  # comment
  if(n == 1) {
    print(true)
  }
}
`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()

	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestError(t *testing.T) {
	code := `let x = 1 `

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	errs := p.Errors()

	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println(e)
		}
		return
	}

	fmt.Println(prog.String())
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

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestSquare(t *testing.T) {
	code := `let x: int = (1,2,3)

x[2] = 3

let y: int = list(1,2,3)

y[[1]] = 1

let zz: string = ("hello|world", "hello|again")
let z: char = strsplit(zz[2], "\\|")[[1]]
`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}

func TestReal(t *testing.T) {
	code := `
#' Type
#'
#' Add type to the roxygen2 documentation.
#'
#' @importFrom roxygen2 roclet roxy_tag_warning block_get_tags roclet_output
#' @importFrom roxygen2 roclet_process roxy_tag_parse rd_section roxy_tag_rd
#'
#' @import roxygen2
#'
#' @export
func roclet_type(): any {
  return roclet("type")
}

#' @export
func (x roxy_tag_type) roxy_tag_parse(): any {
  let parts: char = strsplit(x.raw, ":")
  parts = parts [[1]]

  if(length(parts) != 2){
    roxy_tag_warning("Invalid @type tag, expects <param>: <type> | <type>")
    return
  }

  parts = gsub("\\n|\\t", "", parts)
  let types: char = strsplit(parts[2], "\\|")
  types = types[[1]]

  x.val = list(
    list(
      arg = parts[1] |> trimws(),
      types = types |> trimws()
    )
  )

  return x
}

#' @export
func (x roxy_tag_type) roxy_tag_rd(base_path: char, env: any): any {
  return rd_section("type", x.val)
}

#' @export
func (x rd_section_type) format(...: any): char {
  let types: char = ""
  for (val in x.value) {
    let t: char = paste0(val.types, collapse = ", or ")
    let type: char = paste0("  \\item{", val.arg, "}{", t, "}\n")
    types = paste0(types, type)
  }

  return paste0(
    "\\section{Types}{\n",
    "\\itemize{\n",
    types,
    "}\n",
    "}\n"
  )
}

#' @export
func (x roclet_type) roclet_process(blocks: any, env: any, base_path: char): any {
  let results: any = list()

  for (block in blocks) {
    let tags: any = block_get_tags(block, "type")
    for(tag in tags){
      let t: any = list(
        value = tag.val,
        type = "type",
        file = tag.file
      )
      results = c(results, tag.val)
    }
  }

  return results
}

#' @export
func (x roclet_type) roclet_output(results: any, base_path: char, ...): null {
  .globals.types = results
  return invisible(NULL)
}
`

	l := &lexer.Lexer{
		Input: code,
	}

	l.Run()
	l.Print()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}
