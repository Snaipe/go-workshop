package main

import (
	"fmt"
)

func main() {

	const n = 16

	i := n
	for i != 1 {
		fmt.Print(i, " ")
		i /= 2
	}

	for i := n; i != 1; i /= 2 {
		fmt.Print(i, " ")
	}

	fmt.Println()

	ints := []int{16, 8, 4, 2}
	for index, value := range ints {
		fmt.Printf("[%d]: %d\n", index, value)
	}

	emails := map[string]string{
		"John Doe": "jdoe@acme.org",
		"Snaipe":   "me@snai.pe",
	}
	for key, value := range emails {
		fmt.Printf("%s: %s\n", key, value)
	}

}
