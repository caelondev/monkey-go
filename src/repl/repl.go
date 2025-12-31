package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/caelondev/monkey/src/run"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		run.RunSource(line, out)
	}
}
