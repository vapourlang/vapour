package main

import (
	"log"
	"os"

	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
	"github.com/devOpifex/vapour/transpiler"
	"github.com/devOpifex/vapour/walker"
)

func (v *vapour) transpile(conf Cli) {
	err := v.readDir()

	if err != nil {
		log.Fatal("Failed to read files")
	}

	// lex
	l := &lexer.Lexer{
		Input: string(v.combined),
	}
	l.Run()

	// parse
	p := parser.New(l)
	prog := p.Run()

	// walk tree
	w := walker.New()
	w.Walk(prog)
	if w.HasError() {
		w.Errors().Print()
		return
	}

	// transpile
	trans := transpiler.New()
	trans.Transpile(prog)
	code := trans.GetCode()

	// write
	path := *conf.out + "/vp.R"
	f, err := os.Create(path)

	if err != nil {
		log.Fatal("Failed to create file")
	}

	defer f.Close()

	f.WriteString(code)
}
