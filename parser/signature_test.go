package parser

import (
	"testing"

	"github.com/vapourlang/vapour/ast"
	"github.com/vapourlang/vapour/lexer"
)

func TestFunctionTypeReturn(t *testing.T) {
	l := lexer.NewTest("type math: func(int, int) int\n")
	l.Run()

	p := New(l)
	prog := p.Run()

	for _, s := range prog.Statements {
		tf, ok := s.(*ast.TypeFunction)

		if !ok {
			continue
		}

		if len(tf.Arguments) != 2 {
			t.Fatalf("expected 2 arguments, got %v", len(tf.Arguments))
		}

		if len(tf.Return) != 1 || tf.Return[0].Name != "int" {
			t.Fatalf("expected return type int, got %v", tf.Return)
		}

		return
	}

	t.Fatal("function type not parsed")
}

func TestDecoratorArguments(t *testing.T) {
	codes := []string{
		"@factor(levels = (\"a\", \"b\"))\ntype fct: factor { char }\n",
		"@factor(levels = c(\"a\", \"b\"), ordered = true)\ntype fct: factor { char }\n",
		"@factor(ordered = true, levels = c(\"a\", \"b\"))\ntype fct: factor { char }\n",
		"@environment(parent = emptyenv())\ntype e: environment { x: int }\n",
		"@matrix(nrow = 2, ncol = 2)\ntype m: matrix { int }\n",
	}

	for _, code := range codes {
		l := lexer.NewTest(code)
		l.Run()

		p := New(l)
		p.Run()

		if p.HasError() {
			t.Fatalf("failed to parse:\n%v\n%v", code, p.Errors())
		}
	}
}
