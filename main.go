package main

import (
	"fmt"
	"os"

	"github.com/devOpifex/vapour/cli"
	"github.com/devOpifex/vapour/lsp"
)

func main() {

	v := New()

	args := cli.Cli()

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
	}

	if *args.Version {
		fmt.Printf("v%v\n", v.version)
	}
}
