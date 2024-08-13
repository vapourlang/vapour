package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/devOpifex/vapour/cli"
	"github.com/devOpifex/vapour/lsp"
)

func run(code string) {
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

func (v *vapour) Run(args cli.CLI) {
	if *args.Indir != "" {
		v.transpile(args)
		return
	}

	if *args.Infile != "" {
		v.transpileFile(args)
		return
	}

	if *args.Repl {
		v.repl(os.Stdin, os.Stdout, os.Stderr)
		return
	}

	if *args.LSP {
		lsp.Run(*args.TCP, *args.Port)
		return
	}

	if *args.Version {
		fmt.Printf("v%v\n", v.version)
		return
	}

	fmt.Println("nothing to do, pass at least -repl, -infile, or -indir")
}
