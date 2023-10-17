package main

import (
	"fmt"
)

func main() {
	// Arrays have a fixed-size
	var array [3]int
	array[0] = 1
	array[1] = 2
	array[2] = 3

	fmt.Println("array:", array, "length:", len(array))

	// Does not compile:
	//
	//n := 3
	//var array [n]int

	ints := []int{1, 2, 3}
	fmt.Println("ints:", ints, "length:", len(ints), "capacity:", cap(ints))

	intsAlt := make([]int, 3)
	intsAlt[0] = 1
	intsAlt[1] = 2
	intsAlt[2] = 3
	fmt.Println("ints (alt):", intsAlt)

	// Slicing can be done on both arrays and slices
	fmt.Printf("array[1:2]: %v (type %T, type of array: %T)\n", array[1:2], array[1:2], array)
	fmt.Printf("ints[1:2]: %v (type %T, type of ints: %T)\n", ints[1:2], array[1:2], ints)

	// Slicing does not copy the underlying array
	slice := array[1:2]
	fmt.Println("slice:", slice)
	slice[0] = 4
	fmt.Println("original:", array)

	// Appending is possible, and may allocate a different array
	ints = append(ints, 4)
	fmt.Printf("ints: (@%p) %v length: %d capacity: %d\n", ints, ints, len(ints), cap(ints))
	ints = append(ints, 5, 6)
	fmt.Printf("ints: (@%p) %v length: %d capacity: %d\n", ints, ints, len(ints), cap(ints))
	ints = append(ints, 7)
	fmt.Printf("ints: (@%p) %v length: %d capacity: %d\n", ints, ints, len(ints), cap(ints))
}
