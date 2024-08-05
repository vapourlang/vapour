package main

import (
	"os"

	"github.com/devOpifex/vapour/lsp"
)

func main() {

	v := New()

	args := v.cli()

	if *args.indir != "" {
		v.transpile(args)
		return
	}

	if *args.infile != "" {
		v.transpileFile(args)
		return
	}

	if *args.repl {
		v.repl(os.Stdin, os.Stdout, args)
		return
	}

	if *args.lsp {
		lsp.Run()
	}
}
