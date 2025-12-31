package src

import (
	"fmt"
	"os"

	"github.com/caelondev/monkey/src/repl"
	"github.com/caelondev/monkey/src/run"
)

func Main() {
	args := os.Args

	if len(args) == 1 {
		repl.Start(os.Stdin, os.Stdout)
	} else if len(args) == 2 {
		run.RunFile(args[1])
	} else {
		fmt.Printf("Usage: monkey [filepath]")
		os.Exit(0)
	}
}
