package cli

import (
	"flag"
)

type CLI struct {
	Indir    *string
	Outdir   *string
	LSP      *bool
	TCP      *bool
	Port     *string
	Repl     *bool
	Help     *bool
	Version  *bool
	Check    *bool
	Run      *bool
	Types    *string
	Infile   *string
	Outfile  *string
	Devtools *string
}

func Cli() CLI {
	// inputs
	indir := flag.String("indir", "", "Directory of vapour files to process")
	outdir := flag.String("outdir", "R", "Directory where to place transpiled files from `dir` (defaults to R)")
	infile := flag.String("infile", "", "Vapour file to process")
	outfile := flag.String("outfile", "vapour.R", "Name of R file to where to palce transpiled `infile`. (defaults to vapour.R)")

	// types
	types := flag.String("types", "inst/types.vp", "Path where to generate the type files, only applies if passing a directory with -indir")

	// run type checker
	check := flag.Bool("check-only", false, "Run type checker")

	// lsp
	lsp := flag.Bool("lsp", false, "Run the language server protocol")
	tcp := flag.Bool("tcp", false, "Run the language server protocol on TCP")
	port := flag.String("port", "3000", "Port on which to run the language server protocol, only used if -tcp flag is passed (defaults to 3000)")

	// run
	run := flag.Bool("run-only", false, "Run the transpiled vapour files")

	// repl
	repl := flag.Bool("repl", false, "Run the REPL")

	// version
	version := flag.Bool("version", false, "Retrieve vapour version")

	// devtools
	devtools := flag.String("devtools", "", "Run {devtools} functions after transpilation, accepts `document`, `check`, `install`, separate by comma (e.g.: `document,check`)")

	flag.Parse()

	return CLI{
		Indir:    indir,
		Outdir:   outdir,
		LSP:      lsp,
		TCP:      tcp,
		Port:     port,
		Infile:   infile,
		Outfile:  outfile,
		Repl:     repl,
		Check:    check,
		Run:      run,
		Version:  version,
		Types:    types,
		Devtools: devtools,
	}
}
