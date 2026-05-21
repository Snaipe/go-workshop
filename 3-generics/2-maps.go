package main

import (
	"fmt"
)

func Clone[M ~map[K]V, K comparable, V any](m M) M {
	cloned := make(M, len(m))
	for k, v := range m {
		cloned[k] = v
	}
	return cloned
}

func main() {
	m1 := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}

	fmt.Println(Clone(m1))

	type MyMap map[string]int

	m2 := MyMap{
		"foo": 1,
		"bar": 2,
	}

	fmt.Println(Clone(m2))
}
