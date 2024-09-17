package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/vapourlang/vapour/cli"
	"github.com/vapourlang/vapour/config"
	"github.com/vapourlang/vapour/devtools"
	"github.com/vapourlang/vapour/environment"
	"github.com/vapourlang/vapour/lsp"
	"github.com/vapourlang/vapour/r"
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
	v.config = config.ReadConfig()

	environment.SetLibrary(r.LibPath())

	if *args.Indir != "" {
		ok := v.transpile(args)
		devtools.Run(ok, args)
		return
	}

	if *args.Infile != "" {
		ok := v.transpileFile(args)
		devtools.Run(ok, args)
		return
	}

	if *args.Repl {
		v.repl(os.Stdin, os.Stdout, os.Stderr)
		return
	}

	if *args.LSP {
		lsp.Run(v.config, *args.TCP, *args.Port)
		return
	}

	if *args.Version {
		fmt.Printf("v%v\n", v.version)
		return
	}
}
