package main

import (
	"fmt"
	"github.com/sanemat/go-milligo"
	"github.com/sanemat/go-milligo/codegen"
	"github.com/sanemat/go-milligo/parse"
	"github.com/sanemat/go-milligo/tokenize"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "args is not 2. got=%d", len(os.Args))
		os.Exit(1)
	}

	// Tokenize and parse
	milligo.UserInput = os.Args[1]
	var err error
	milligo.Tk, err = tokenize.Tokenize()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
	node := parse.Expr()
	codegen.Codegen(node)
}
