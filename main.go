package main

import (
	"github.com/devOpifex/vapour/cli"
)

func main() {
	v := New()
	args := cli.Cli()
	v.Run(args)
}
