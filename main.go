package main

import (
	"fmt"
	"os"

	"github.com/zachklingbeil/challenge/fx"
)

// main expects a single argument: the path to the input file.
// If the argument is missing, it prints usage instructions and exits.
// Otherwise, output uptime results.
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Need input to output.\nelectric <input file>")
		return
	}

	era := fx.Electric()   // initialize new Era
	era.Input(os.Args[1])  // parse
	era.Fx()               // calculate
	era.Output(era.Uptime) // results
}
