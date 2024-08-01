package main

import (
	"fmt"
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
		log.Fatal("Failed to read vapour files")
	}

	// lex
	l := lexer.New(v.files)
	l.Run()

	if l.HasError() {
		l.Errors.Print()
		return
	}

	// parse
	p := parser.New(l)
	prog := p.Run()

	if p.HasError() {
		for _, e := range p.Errors() {
			fmt.Println(e)
		}
		return
	}

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

	_, err = f.WriteString(code)

	if err != nil {
		log.Fatal("Failed to write to file")
	}
}
