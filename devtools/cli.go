package devtools

import (
	"fmt"
	"strings"

	"github.com/vapourlang/vapour/cli"
	"github.com/vapourlang/vapour/r"
)

func Run(valid bool, args cli.CLI) {
	if !valid {
		return
	}

	if *args.Devtools == "" {
		return
	}

	commands := strings.Split(*args.Devtools, ",")

	for _, c := range commands {
		cmd := "devtools::" + c + "()"
		out, err := r.Callr(cmd)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Println(string(out[:]))
	}
}
