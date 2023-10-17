package main

import (
	"fmt"
)

func main() {

	m := map[string]int{
		"a": 1,
	}

	m["b"] = 2

	delete(m, "b")

	fmt.Println(m, len(m))

}
