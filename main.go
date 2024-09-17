package main

import (
	"github.com/vapourlang/vapour/cli"
)

func main() {
	v := New()
	args := cli.Cli()
	v.Run(args)
}
