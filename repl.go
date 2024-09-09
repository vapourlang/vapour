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
)

const PROMPT = "> "

func (v *vapour) replIntro() string {
	return "Vapour " + v.version + "\n"
}

func (v *vapour) repl(in io.Reader, out io.Writer, er io.Writer) {
	cmd := exec.Command(
		"R",
		"-s",
		"-e",
		`f <- file("stdin")
    open(f)
		while(length(line <- readLines(f, n = 1)) > 0) {
			cat(eval(parse(text = line)))
		}`,
	)

	//cmd.Stdout = out
	cmd.Stderr = er

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

	fmt.Fprint(out, v.replIntro())

	err = cmd.Start()

	if err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Fprint(out, PROMPT)

		scanned := scanner.Scan()

		if !scanned {
			continue
		}

		line := scanner.Text()

		// lex
		l := lexer.NewCode("repl", line+"\n")
		l.Run()

		if l.HasError() {
			fmt.Fprintln(out, l.Errors().String())
			continue
		}

		// parse
		p := parser.New(l)
		prog := p.Run()

		if p.HasError() {
			for _, e := range p.Errors() {
				fmt.Fprintln(out, e.Message)
			}
			continue
		}

		trans := transpiler.New()
		trans.Transpile(prog)
		code := trans.GetCode()

		_, err := io.WriteString(cmdIn, code)

		if err != nil {
			fmt.Fprint(out, err.Error())
			continue
		}

		res := make([]byte, 1024)
		_, err = cmdOut.Read(res)

		if err != nil {
			fmt.Fprint(out, err.Error())
			continue
		}

		fmt.Fprint(out, string(res)+"\n")
	}
}
