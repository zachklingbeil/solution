package main

import (
	"fmt"
	"os"
	"slices"
)

func Output(m map[uint32]int) {
	keys := make([]uint32, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		fmt.Printf("%d %d\n", k, m[k])
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR")
		return
	}
	NewElectric(os.Args[1])
}
