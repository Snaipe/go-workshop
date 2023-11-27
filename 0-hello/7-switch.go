package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {

	switch len(os.Args) {
	case 1:
		// no break statement necessary
	case 0:
		fallthrough // falls through to case under it
	default:
		fmt.Println("usage: go run 4-if.go <number>")
		os.Exit(2)
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// The value in switch statements can be omitted, and this behaves like
	// an if-elseif-else chain.

	switch {
	case n <= 0:
		fmt.Printf("%d < 0\n", n)
	case n < 100:
		fmt.Printf("%d < 100\n", n)
	default:
		fmt.Printf("%d >= 100\n", n)
	}

}

