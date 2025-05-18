package main

import (
	"fmt"
	"os"

	"github.com/zachklingbeil/solution/fx"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("outputs require inputs")
		fmt.Println("\nmissing path/to/<input>.txt")
		return
	}

	era := fx.Electric()
	era.Input(os.Args[1])
	era.Fx()
	era.Output()
}
