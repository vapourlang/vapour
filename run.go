package main

import (
	"fmt"
	"log"
	"os/exec"
)

func run(code string) {
	cmd := exec.Command(
		"R",
		"--no-save",
		"--slave",
		"-e",
		code,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal("Failed to run")
	}

	fmt.Println(string(output))
}
