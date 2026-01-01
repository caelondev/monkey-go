package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/caelondev/monkey/src/object"
	"github.com/caelondev/monkey/src/run"
	"github.com/jwalton/gchalk"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	var allLines []string

	for {
		fmt.Printf(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		allLines = append(allLines, line)

		result := run.RunSource(line, out)

		if result != nil {
			if result.Type() == object.ERROR_OBJECT {
				err := result.(*object.Error)
				lineColumn := gchalk.WithBold().Red("Runtime::Error")
				message := gchalk.Red(" -> " + err.Message + "\n")
				io.WriteString(out, lineColumn+message)
			} else {
				io.WriteString(out, result.Inspect())
				io.WriteString(out, "\n")
			}
		}
	}
}
