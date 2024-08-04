package main

import (
	"os"
)

func main() {

	v := &vapour{
		name:    "vapour",
		version: "0.0.1",
	}

	args := v.cli()

	if *args.dir != "" {
		v.root = args.dir
		v.transpile(args)
		return
	}

	if *args.run != "" {
		v.run(args)
		return
	}

	if *args.repl {
		v.repl(os.Stdin, os.Stdout, args)
		return
	}

	if *args.help {
		v.help()
		return
	}

	if *args.lsp {
		v.lspRun()
	}
}
