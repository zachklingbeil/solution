package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/zachklingbeil/challenge/fx"
)

//go:embed input1.txt
var input1 string

//go:embed input2.txt
var input2 string

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Need input to output.\nelectric <input file>")
		return
	}

	var source string
	switch os.Args[1] {
	case "1":
		source = "input1.txt"
	case "2":
		source = "input2.txt"
	default:
		source = os.Args[1]
	}

	era := fx.Electric()
	era.Input(source)
	era.Fx()
	era.Output()
}
