package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/caelondev/monkey/src/evaluation"
	"github.com/caelondev/monkey/src/lexer"
	"github.com/caelondev/monkey/src/parser"
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

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			io.WriteString(out, "An error occured whilst parsing:\n")
			printParserErrors(out, p.Errors())
			io.WriteString(out, "\n")
			continue
		}

		evaluator := evaluation.New(line)

		result := evaluator.Evaluate(program)

		if result != nil {
			io.WriteString(out, result.Inspect())
			io.WriteString(out, "\n")
		}

		// Print AST ---
		// io.WriteString(out, program.String())
		// io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
