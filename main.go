package main

import (
	"fmt"
	"os"

	"github.com/zachklingbeil/electric/fx"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Need input to output.\nelectric <input file>")
		return
	}

	era := fx.Electric()
	era.Input(os.Args[1])
	era.Fx()
	era.Output(era.Uptime)
}
