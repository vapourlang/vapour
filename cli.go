package main

import (
	"flag"
)

type CLI struct {
	indir   *string
	outdir  *string
	lsp     *bool
	tcp     *bool
	port    *string
	repl    *bool
	help    *bool
	version *bool
	check   *bool
	run     *bool
	infile  *string
	outfile *string
}

func (v *vapour) cli() CLI {
	// inputs
	indir := flag.String("indir", "", "Directory of vapour files to process")
	outdir := flag.String("outdir", "R", "Directory where to place transpiled files from `dir`")
	infile := flag.String("infile", "", "Vapour file to process")
	outfile := flag.String("outfile", "vapour.R", "Name of R file to where to palce transpiled `infile`")

	// run type checker
	check := flag.Bool("check-only", false, "Run type checker")

	// lsp
	lsp := flag.Bool("lsp", false, "Run the language server protocol")
	tcp := flag.Bool("tcp", false, "Run the language server protocol on TCP")
	port := flag.String("port", "5007", "Port on which to run the language server protocol")

	// run
	run := flag.Bool("run-only", false, "Run the transpiled vapour files")

	// repl
	repl := flag.Bool("repl", false, "Run the REPL")

	// help
	help := flag.Bool("help", false, "Get help on commands")

	// version
	version := flag.Bool("version", false, "Retrieve vapour version")

	flag.Parse()

	return CLI{
		indir:   indir,
		outdir:  outdir,
		lsp:     lsp,
		tcp:     tcp,
		port:    port,
		infile:  infile,
		outfile: outfile,
		repl:    repl,
		help:    help,
		check:   check,
		run:     run,
		version: version,
	}
}
