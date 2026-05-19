package main

import (
	"fmt"
	"strings"
)

// Fibonacci returns the nth element in the fibonacci sequence.
func Fibonacci(n int) int {
	if n == 0 {
		return n
	}

	a, b := 0, 1
	for i := 1; i < n; i++ {
		a, b = b, a+b
	}
	return b
}

// ParseInt parses a base-10 integer from `in`.
func ParseInt(in string) (int, error) {
	sign := 1
	if strings.HasPrefix(in, "-") {
		sign = -1
		in = in[1:]
	}

	var result int
	for i, c := range in {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("ParseInt %q: at index %d: invalid rune %q", in, i, c)
		}
		result = result*10 + int(c-'0')
	}
	return sign * result, nil
}

func main() {

	fmt.Print("fibonacci: ")
	for i := 0; i < 10; i++ {
		fmt.Printf("%d ", Fibonacci(i))
	}
	fmt.Println()

	fmt.Println(ParseInt("1234"))
	fmt.Println(ParseInt("-1234"))
	fmt.Println(ParseInt("1234NotANumber"))

}
