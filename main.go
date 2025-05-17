package main

import (
	"fmt"
	"os"

	v2 "github.com/zachklingbeil/challenge/v2"
)

// main expects a single argument: the path to the input file.
// If the argument is missing, it prints usage instructions and exits.
// Otherwise, output uptime results.
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Need input to output.\nelectric <input file>")
		return
	}

	era := v2.NewFx() // initialize new Era
	era.ParseReportsFromFile(os.Args[1])
	era.CalculateUptime()
	era.PrintStationUptime()
}
