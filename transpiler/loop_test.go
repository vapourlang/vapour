package transpiler

import (
	"testing"

	"github.com/vapourlang/vapour/lexer"
	"github.com/vapourlang/vapour/parser"
)

func TestLoopControl(t *testing.T) {
	code := `for(let i: int in 1..10) {
  if(i == 2) {
    next
  }
  if(i == 5) {
    break
  }
}

repeat {
  break
}
`

	l := lexer.NewTest(code)
	l.Run()
	p := parser.New(l)
	prog := p.Run()

	if p.HasError() {
		t.Fatalf("failed to parse: %v", p.Errors())
	}

	trans := New()
	trans.Transpile(prog)

	expected := `for(i in 1:10
) {
if(i==2
){
next}
if(i==5
){
break}
}
repeat {
break}
`

	trans.testOutput(t, expected)
}
