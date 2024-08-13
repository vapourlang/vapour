package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/parser"
	"github.com/devOpifex/vapour/transpiler"
	"github.com/devOpifex/vapour/walker"
)

const PROMPT = "> "

func (v *vapour) replIntro() string {
	return "Vapour " + v.version + "\n"
}

func (v *vapour) repl(in io.Reader, out io.Writer, conf CLI) {
	cmd := exec.Command(
		"R",
	)

	cmdIn, err := cmd.StdinPipe()

	if err != nil {
		log.Fatal(err)
	}

	defer cmdIn.Close()

	cmdOut, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal(err)
	}

	defer cmdOut.Close()

	scanner := bufio.NewScanner(in)

	fmt.Fprintf(out, v.replIntro())

	cmd.Run()
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		// lex
		l := lexer.NewCode("repl", line+"\n")
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
				io.WriteString(out, e.Message)
			}
			return
		}

		// walk tree
		w := walker.New()
		w.Walk(prog)
		if w.HasError() {
			io.WriteString(out, w.Errors().String())
			return
		}

		// transpile
		trans := transpiler.New()
		trans.Transpile(prog)
		code := trans.GetCode()

		io.WriteString(cmdIn, code)

		res := make([]byte, 1024)
		_, err = cmdOut.Read(res)

		if err != nil {
			io.WriteString(out, err.Error())
			continue
		}

		io.WriteString(out, string(res))
	}
}
