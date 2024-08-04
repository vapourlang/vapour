package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
	"github.com/devOpifex/vapour/transpiler"
	"github.com/devOpifex/vapour/walker"
)

const PROMPT = "> "

func (v *vapour) repl(conf Cli) {
	cmd := exec.Command(
		"R",
	)

	stdin, err := cmd.StdinPipe()

	if err != nil {
		log.Fatal(err)
	}

	start(os.Stdin, os.Stdout, stdin)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(out)
}

func start(in io.Reader, out io.Writer, stdin io.WriteCloser) {
	scanner := bufio.NewScanner(in)
	defer stdin.Close()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		// lex
		l := lexer.NewCode("repl", line)
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

		io.WriteString(stdin, code)
	}
}
