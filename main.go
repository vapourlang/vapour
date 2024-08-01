package main

import "fmt"

func main() {

	v := &vapour{
		name: "vapour",
	}

	args := v.cli()

	if *args.dir != "" {
		v.root = args.dir
		v.transpile(args)
		return
	}

	if *args.lsp {
		v.lspInit()
		v.lspRun()
		return
	}

	fmt.Println("Nothing to do")
}
