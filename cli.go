package main

import "flag"

type Cli struct {
	dir *string
	out *string
	lsp *bool
	run *string
}

func (v *vapour) cli() Cli {
	// transpile
	dir := flag.String("dir", "", "Directory of .vp files to transpile vapour to R")
	out := flag.String("out", "R", "Directory where to place transpiled files")

	// lsp
	lsp := flag.Bool("lsp", false, "Run the language server protocol")

	// repl
	run := flag.String("run", "", "Run vapour code")

	flag.Parse()

	return Cli{
		dir: dir,
		out: out,
		lsp: lsp,
		run: run,
	}
}
