package main

import (
	"fmt"
)

func main() {

	v1 := 1

	ptr := &v1

	fmt.Println(ptr, v1)

	*ptr = 2

	fmt.Println(v1)

	ints := []int{1, 2, 3}

	fmt.Println(&ints[0], &ints[1], &ints[2])

}

func Swap(a, b *int) {
	*a, *b = *b, *a
}
