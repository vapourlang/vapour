package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/vapourlang/vapour/lexer"
	"github.com/vapourlang/vapour/parser"
	"github.com/vapourlang/vapour/transpiler"
)

const PROMPT = "> "

// replDone marks the end of R's output for one evaluation
const replDone = "\x01VAPOUR_DONE\x01"

// replR reads one quoted R string per line from stdin, evaluates the
// code it contains, prints visible results as the R console would,
// then prints replDone so the REPL knows the evaluation is complete.
const replR = `f <- file("stdin")
open(f)
while(length(line <- readLines(f, n = 1)) > 0) {
  tryCatch({
    code <- eval(parse(text = line))
    res <- withVisible(eval(parse(text = code), envir = globalenv()))
    if(res$visible) print(res$value)
  }, error = function(e) {
    cat("Error:", conditionMessage(e), "\n")
  })
  cat("` + replDone + `\n")
}`

func (v *vapour) replIntro() string {
	return "Vapour " + v.version + "\n"
}

func (v *vapour) repl(in io.Reader, out io.Writer, er io.Writer) {
	cmd := exec.Command("R", "-s", "-e", replR)
	cmd.Stderr = er

	cmdIn, err := cmd.StdinPipe()

	if err != nil {
		log.Fatal(err)
	}

	cmdOut, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Start()

	if err != nil {
		log.Fatal(err)
	}

	// closing stdin ends the R loop, then wait for R to exit
	defer cmd.Wait()
	defer cmdIn.Close()

	scanner := bufio.NewScanner(in)
	reader := bufio.NewReader(cmdOut)

	fmt.Fprint(out, v.replIntro())

	for {
		fmt.Fprint(out, PROMPT)

		// end of input, e.g.: Ctrl+D
		if !scanner.Scan() {
			fmt.Fprintln(out)
			return
		}

		line := scanner.Text()

		if strings.TrimSpace(line) == "" {
			continue
		}

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

		// send the code as a single quoted line, the transpiled
		// code may span multiple lines
		_, err := io.WriteString(cmdIn, strconv.Quote(code)+"\n")

		if err != nil {
			fmt.Fprintln(out, err.Error())
			return
		}

		// print R's output until the end marker
		for {
			res, err := reader.ReadString('\n')

			if err != nil {
				fmt.Fprintln(out, err.Error())
				return
			}

			if strings.TrimRight(res, "\n") == replDone {
				break
			}

			fmt.Fprint(out, res)
		}
	}
}
