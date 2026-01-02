package run

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/caelondev/monkey/src/evaluation"
	"github.com/caelondev/monkey/src/lexer"
	"github.com/caelondev/monkey/src/object"
	"github.com/caelondev/monkey/src/parser"
	"github.com/jwalton/gchalk"
)

var ENVIRONMENT *object.Environment = object.NewEnvironment(nil)

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

	source := string(byte)
	result := RunSource(source, os.Stdout)

	if result != nil && result.Type() == object.ERROR_OBJECT {
		formatFileError(result.(*object.Error), source, os.Stdout)
	}
}

func RunSource(source string, out io.Writer) object.Object {
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		io.WriteString(out, "An error occured whilst parsing:\n")
		printParserErrors(out, p.Errors())
		io.WriteString(out, "\n")
		return nil
	}

	evaluator := evaluation.New()

	result := evaluator.Evaluate(program, ENVIRONMENT)

	return result
}

func formatFileError(err *object.Error, source string, out io.Writer) {
	lines := strings.Split(source, "\n")

	lineColumn := gchalk.WithBold().Red(fmt.Sprintf("[Ln %d:%d] Runtime::Error", err.Line, err.Column))
	message := gchalk.Red(" -> " + err.Message)

	snippet := "\n\n"

	if int(err.Line) > 0 && int(err.Line) <= len(lines) {
		sourceLine := lines[err.Line-1]
		lineNumStr := fmt.Sprintf("Ln %d:%d", err.Line, err.Column)

		snippet += gchalk.WithBold().White(" Error caused by:\n")
		snippet += gchalk.Cyan(fmt.Sprintf("    %s | ", lineNumStr))
		snippet += gchalk.White(sourceLine + "\n")

		padding := strings.Repeat(" ", len(lineNumStr))
		pointer := strings.Repeat(" ", int(err.Column)-1) + "^"
		snippet += gchalk.Cyan(fmt.Sprintf("    %s | ", padding))
		snippet += gchalk.BrightRed(pointer + "\n")
	} else {
		snippet += gchalk.WithBold().White(" Error caused by:\n")
		snippet += gchalk.Cyan(fmt.Sprintf("\t%d:%d | ", err.Line, err.Column))
		snippet += gchalk.White(err.NodeStr + "\n")
	}

	if err.Hint != "" {
		snippet += "\n"
		snippet += gchalk.Cyan(" Hint: ")
		snippet += gchalk.White(err.Hint + "\n")
	}

	io.WriteString(out, lineColumn+message+snippet)
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
