package transpiler

import (
	"fmt"
	"testing"

	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
)

func TestBasic(t *testing.T) {
	code := `let x: int | num = 1  `

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
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

	fmt.Println(trans.GetCode())
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

	fmt.Println(trans.GetCode())
}

func TestString(t *testing.T) {
	code := `let x: string = "a \"string\""
let y: string = 'single quotes'`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
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

	fmt.Println(trans.GetCode())
}

func TestMethod(t *testing.T) {
	code := `func (o obj) add(n: int): char {
  return "hello"
}`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestTypeDeclaration(t *testing.T) {
	code := `type userId: int

type obj: struct {
  int | string,
  name: string,
  id: int
} `

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestIdent(t *testing.T) {
	code := `let x: int = (1,2,3)

x[1, 2] = 15

x[[3]] = 15

df.x = 23

print(x) `

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestTElipsis(t *testing.T) {
	code := `func foo(...: any) char {
  paste0(..., collapse = ", ")
}  `

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
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

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
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

	fmt.Println(trans.GetCode())
}

func TestFor(t *testing.T) {
	code := `for(let i:int = 1 in 1..nrow(df)) {
  print(i)
}

func foo(...: int): int {
  sum(...)
}

let x: int = (1, 20, 23) `

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
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

	fmt.Println(trans.GetCode())
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

	fmt.Println(trans.GetCode())
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

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
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

	fmt.Println(trans.GetCode())
}

func TestSquare(t *testing.T) {
	code := `let x: int = (1,2,3)

x[2] = 3

let y: int = list(1,2,3)

y[[1]] = 1

let zz: string = ("hello|world", "hello|again")
let z: char = strsplit(zz[2], "\\|")[[1]]
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}

func TestReal(t *testing.T) {
	code := `type Config: list {
  types: string,
  yields: string
}

func write_config() {
  let config: Config = Config(
    types = globals.types,
    yields = globals.yields
  )

  print(config)
}

#' @export
func roclet_type(): any {
  return roclet("type")
}

#' @export
func (x roxy_tag_type) roxy_tag_parse(): any {
  let parts: char = strsplit(x.raw, ":")
  parts = parts[[2]]

  if(length(parts) != 2){
    roxy_tag_warning("Invalid @type tag, expects <param>: <type> | <type>")
    return null
  }

  parts = gsub("\\n|\\t", "", parts)
  let types: char = strsplit(parts[2], "\\|")[[1]]

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

func (x rd_section_type) format(...: any): char {
  let types: char = ""
  for (val in x.value) {
    let t: char = paste0(val.types, collapse = ", or ")
    t = paste0("  \\item{", val.arg, "}{", t, "}\n")
    types = paste0(types, t)
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
        cat = "type",
        file = tag.file
      )
      results = c(results, tag.val)
    }
  }

  return results
}

#' @export
func (x roclet_type) roclet_output(results: any, base_path: char, ...): null {
  globals.types = results
  return invisible(NULL)
}

type roxy: list {
  yield: string
}

#' Yield
#'
#' Add yield to the roxygen2 documentation.
#'
#' @importFrom roxygen2 roclet roxy_tag_warning block_get_tags roclet_output
#' @importFrom roxygen2 roclet_process roxy_tag_parse rd_section roxy_tag_rd
#'
#' @import roxygen2
#'
#' @export
func roclet_yield() {
  return roclet("yield")
}

#' @export
func (x rody_tag_yield) roxy_tag_parse(x: any): any {
  let raw: char = gsub("\\n|\\t", "", x.raw)
  let yields: char = strsplit(raw, "\\|")[[1]]

  x.val: roxy = roxy(
    yield = yields |> trimws()
  )

  return x
}

#' @export
func (x roxy_tag_yield) rody_tag_rd(base_path: char, env: any): any {
  rd_section("yield", x.val)
}

#' @export
func(x rd_section_yield) format(...: any): char {
  let yield: char = paste0(x.value.yield, collapse = ", or ")
  return paste0(
    "\\yield{", yield, "}\n"
  )
}

#' @export
func (x roclet_yield) roclet_process(blocks: any, env: any, base_path: char): list {
  let results: list = list()

  for (block in blocks) {
    let tags: any = block_get_tags(block, "yield")
    class(tags) = "list"
    results = append(results, list(x))
  }

  return results
}

#' @export
func (x roclet_yield) roclet_output(results: any, base_path: char, ...: any): null {
  globals.yields = results
  return invisible(null)
}
`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	trans := New()
	trans.Transpile(prog)

	fmt.Println(trans.GetCode())
}
