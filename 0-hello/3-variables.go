package main

import (
	"fmt"
	"strings"
)

func main() {

	v1 := 1

	// Equivalent to v1
	var v2 int = 1
	var v3 = 1
	fmt.Println("v1:", v1, "v2:", v2, "v3:", v3)

	// Zero values
	var (
		s string // = ""
		i int    // = 0
	)
	fmt.Println("s:", s, "i:", i)

	// Multiple values
	s, i = "", 0

	before, after, found := strings.Cut("first,second", ",")
	fmt.Println("before:", before, "after:", after, "found:", found)

}

var (
	ExportedGlobal = 1

	unexportedGlobal = 2
)
