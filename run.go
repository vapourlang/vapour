package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
	"github.com/devOpifex/vapour/transpiler"
	"github.com/devOpifex/vapour/walker"
)

func (v *vapour) run(conf Cli) {
	content, err := os.ReadFile(*conf.run)

	if err != nil {
		log.Fatal("Could not read vapour file")
	}

	// lex
	l := lexer.NewCode(*conf.run, string(content))
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

	cmd := exec.Command(
		"R",
		"--no-save",
		"--slave",
		"-e",
		code,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal("Failed to run")
	}

	fmt.Println(string(output))
}
