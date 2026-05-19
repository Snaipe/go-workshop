package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) <= 1 {
		fmt.Println("usage: go run 4-if.go <number>")
		os.Exit(2)
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if n <= 0 {
		fmt.Printf("%d < 0\n", n)
	} else if n < 100 {
		fmt.Printf("%d < 100\n", n)
	} else {
		fmt.Printf("%d >= 100\n", n)
	}

}
