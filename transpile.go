package main

import (
	"fmt"
	"log"
	"os"

	"github.com/devOpifex/vapour/cli"
	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
	"github.com/devOpifex/vapour/transpiler"
	"github.com/devOpifex/vapour/walker"
)

func (v *vapour) transpile(conf cli.CLI) bool {
	v.root = conf.Indir
	err := v.readDir()

	if err != nil {
		log.Fatal("Failed to read vapour files")
	}

	// lex
	l := lexer.New(v.files)
	l.Run()

	if l.HasError() {
		l.Errors().Print()
		transpileFailed()
		return false
	}

	// parse
	p := parser.New(l)
	prog := p.Run()

	if p.HasError() {
		p.Errors().Print()
		transpileFailed()
		return false
	}

	// walk tree
	w := walker.New()
	w.Walk(prog)

	if w.HasDiagnostic() {
		w.Errors().Print()
	}

	if w.HasError() {
		transpileFailed()
		return false
	}

	if *conf.Check {
		return false
	}

	// transpile
	trans := transpiler.New()
	trans.Transpile(prog)
	code := trans.GetCode()

	transpileSuccessful()

	if *conf.Run {
		run(code)
		return false
	}

	code = addHeader(code)

	// write
	path := *conf.Outdir + "/" + *conf.Outfile
	f, err := os.Create(path)

	if err != nil {
		log.Fatalf("Failed to create output file: %v", err.Error())
	}

	defer f.Close()

	_, err = f.WriteString(code)

	if err != nil {
		log.Fatalf("Failed to write output file: %v", err.Error())
	}

	// we only generate types if it's an R package
	if *conf.Outdir != "R" {
		return false
	}

	// write types
	lines := trans.Env().GenerateTypes().String()
	f, err = os.Create(*conf.Types)

	if err != nil {
		log.Fatalf("Failed to create type file: %v", err.Error())
	}

	defer f.Close()

	_, err = f.WriteString(lines)

	if err != nil {
		log.Fatalf("Failed to write to types file: %v", err.Error())
	}

	return true
}

func (v *vapour) transpileFile(conf cli.CLI) bool {
	content, err := os.ReadFile(*conf.Infile)

	if err != nil {
		log.Fatal("Could not read vapour file")
	}

	// lex
	l := lexer.NewCode(*conf.Infile, string(content))
	l.Run()

	if l.HasError() {
		l.Errors().Print()
		transpileFailed()
		return false
	}

	// parse
	p := parser.New(l)
	prog := p.Run()

	if p.HasError() {
		p.Errors().Print()
		transpileFailed()
		return false
	}

	// walk tree
	w := walker.New()
	w.Walk(prog)
	if w.HasDiagnostic() {
		w.Errors().Print()
		transpileFailed()
	}

	if w.HasError() {
		return false
	}

	if *conf.Check {
		return false
	}

	// transpile
	trans := transpiler.New()
	trans.Transpile(prog)
	code := trans.GetCode()

	transpileSuccessful()

	if *conf.Run {
		run(code)
		return false
	}

	code = addHeader(code)

	// write
	f, err := os.Create(*conf.Outfile)

	if err != nil {
		log.Fatal("Failed to create output file")
	}

	defer f.Close()

	_, err = f.WriteString(code)

	if err != nil {
		log.Fatal("Failed to write to output file")
	}

	return true
}

func transpileSuccessful() {
	fmt.Println(cli.Green + "✓" + cli.Reset + " files successfully transpiled!")
}

func transpileFailed() {
	fmt.Println(cli.Red + "x" + cli.Reset + " failed to transpile files!")
}
