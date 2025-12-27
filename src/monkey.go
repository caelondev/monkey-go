package src

import (
	"fmt"

	"github.com/caelondev/monkey/src/lexer"
	"github.com/caelondev/monkey/src/parser"
	"github.com/sanity-io/litter"
)

func Main() {
	src := `
	var = = 10;
	var y = 15;
	var foobar = 192882;
	`

	l := lexer.New(src)
	p := parser.New(l)

	program := p.ParseProgram()

	fmt.Printf("Got %d errors\n", len(p.Errors()))

	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			fmt.Printf("Error: %s\n", err)
		}
	} else {
		litter.Dump(program) // Only dump if no errors
	}
}
