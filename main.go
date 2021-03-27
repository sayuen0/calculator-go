package main

import (
	lg "github.com/sayuen0/calculator-go/lex"
	"os"
)

func main() {
	var lex lg.Lex
	lex.Init(os.Stdin)
	for {
		if lg.TopLevel(&lex) {
			break
		}
	}
}
