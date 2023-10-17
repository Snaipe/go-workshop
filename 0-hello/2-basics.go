package main

import (
	"fmt"
	"math"
)

func main() {

	// Strings
	fmt.Println("string =", "foo"+"bar")

	// Integers
	fmt.Println("int =", 42)
	fmt.Println("sum =", 1+2+3+4+5+6+7+8+9)

	// IEEE Floats
	fmt.Println("third =", 1./3.)
	fmt.Println("φ =", (1+math.Sqrt(5))/2)

}
