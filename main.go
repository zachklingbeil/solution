package main

import (
	"fmt"
	"os"

	"github.com/zachklingbeil/solution/fx"
)

// main is the entry point of the program. It expects exactly one command-line argument: the path to an input file.
// The program reads the input file, processes charger uptime data, and prints the results.
func main() {
	if len(os.Args) != 2 {
		// Print usage instructions if the input file is not provided
		fmt.Println("outputs require inputs. try again using:")
		fmt.Println("./solution input1.txt\nor\n./solution input2.txt")
		return
	}

	// Create a new Era instance to hold station and charger data
	era := fx.Electric()
	// Read and parse the input file
	era.Input(os.Args[1])
	// Calculate uptime percentages for each station
	era.Fx()
	// Output the results to standard output
	era.Output()
}
