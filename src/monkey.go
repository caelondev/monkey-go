package src

import (
	"os"

	"github.com/caelondev/monkey/src/repl"
)

func Main() {
	repl.Start(os.Stdin, os.Stdout)
}
