package run

import (
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"github.com/caelondev/monkey/src/evaluation"
	"github.com/caelondev/monkey/src/lexer"
	"github.com/caelondev/monkey/src/object"
	"github.com/caelondev/monkey/src/parser"
)

func RunFile(filepath string) {
	byte, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("An error occurred whilst trying to read file:\n%s", err.Error())
		os.Exit(1)
	}

	if !utf8.Valid(byte) {
		fmt.Printf("Cannot read non-UTF8 file\n")
		os.Exit(2)
		return
	}

	RunSource(string(byte), os.Stdout)
}

func RunSource(source string, out io.Writer) {
	l := lexer.New(source)
	p := parser.New(l)
	env := object.NewEnvironment()

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		io.WriteString(out, "An error occured whilst parsing:\n")
		printParserErrors(out, p.Errors())
		io.WriteString(out, "\n")
		return
	}

	evaluator := evaluation.New(source)

	result := evaluator.Evaluate(program, env)

	if result != nil {
		io.WriteString(out, result.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
