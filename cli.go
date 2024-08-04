package main

import (
	"flag"
	"io"
	"os"
)

type Cli struct {
	dir     *string
	out     *string
	lsp     *bool
	repl    *bool
	help    *bool
	version *bool
	run     *string
}

func (v *vapour) cli() Cli {
	// transpile
	dir := flag.String("dir", "", "Directory of .vp files to transpile vapour to R")
	out := flag.String("out", "R", "Directory where to place transpiled files")

	// lsp
	lsp := flag.Bool("lsp", false, "Run the language server protocol")

	// run
	run := flag.String("run", "", "Run vapour code")

	// repl
	repl := flag.Bool("repl", false, "Run the REPL")

	// help
	help := flag.Bool("help", false, "Get Help")

	// version
	version := flag.Bool("version", false, "Get vapour version")

	flag.Parse()

	return Cli{
		dir:     dir,
		out:     out,
		lsp:     lsp,
		run:     run,
		repl:    repl,
		help:    help,
		version: version,
	}
}

func (v *vapour) help() {
	io.WriteString(
		os.Stdout,
		"Vapour v"+v.version+", commands:\n"+"-dir -out\n-lsp\n-help\n-version\n",
	)
}
