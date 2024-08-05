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
	v.root = conf.indir
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

	if *conf.check && !w.HasError() {
		fmt.Println("No issues found!")
		return
	}

	if *conf.check {
		return
	}

	// transpile
	trans := transpiler.New()
	trans.Transpile(prog)
	code := trans.GetCode()

	if *conf.run {
		run(code)
		return
	}

	code = addHeader(code)

	// write
	path := *conf.outdir + "/" + *conf.outfile
	f, err := os.Create(path)

	if err != nil {
		log.Fatal("Failed to create file")
	}

	defer f.Close()

	_, err = f.WriteString(code)

	if err != nil {
		log.Fatal("Failed to write to output file")
	}
}

func (v *vapour) transpileFile(conf Cli) {
	content, err := os.ReadFile(*conf.infile)

	if err != nil {
		log.Fatal("Could not read vapour file")
	}

	// lex
	l := lexer.NewCode(*conf.infile, string(content))
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

	if *conf.check && !w.HasError() {
		fmt.Println("No issues found!")
		return
	}

	if *conf.check {
		return
	}

	// transpile
	trans := transpiler.New()
	trans.Transpile(prog)
	code := trans.GetCode()

	if *conf.run {
		run(code)
		return
	}

	code = addHeader(code)

	// write
	f, err := os.Create(*conf.outfile)

	if err != nil {
		log.Fatal("Failed to create output file")
	}

	defer f.Close()

	_, err = f.WriteString(code)

	if err != nil {
		log.Fatal("Failed to write to output file")
	}
}
