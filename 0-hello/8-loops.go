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
	fmt.Println()

	for i := n; i != 1; i /= 2 {
		fmt.Print(i, " ")
	}
	fmt.Println()

	// Iterating on slices

	ints := []int{16, 8, 4, 2}
	for index, value := range ints {
		fmt.Printf("[%d]: %d\n", index, value)
	}
	fmt.Println()

	// Iterating on strings

	s := "Hello, 世界"

	// Generally incorrect
	for i := 0; i < len(s); i++ {
		fmt.Printf("[%d]: %q\n", i, s[i])
	}
	fmt.Println()

	// Correct
	for i, r := range s {
		fmt.Printf("[%d]: %q\n", i, r)
	}
	fmt.Println()

	// Iterating on maps
	//
	// NOTE: The iteration order is randomized

	emails := map[string]string{
		"John Doe": "jdoe@acme.org",
		"Snaipe":   "me@snai.pe",
	}
	for key, value := range emails {
		fmt.Printf("%s: %s\n", key, value)
	}

}
