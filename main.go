package main

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

	if *args.run != "" {
		v.run(args)
		return
	}

	if *args.repl {
		v.repl(args)
		return
	}

	if *args.lsp {
		v.lspInit()
		v.lspRun()
	}
}
