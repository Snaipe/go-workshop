package main

func main() {

	vals := []string{"a", "b", "c"}

	// Prints 0 a, 1 b, 2 c
	for i, v := range vals {
		fmt.Println(i, v)
	}

	// Prints 0 a, 1 b, 2 c
	for i, v := range vals {
		go fmt.Println(i, v)
	}

	// Prints 2 c, 2 c, 2 c
	for i, v := range vals {
		go func() {
			fmt.Println(i, v)
		}()
	}

	// Prints 0 a, 1 b, 2 c
	for i, v := range vals {
		i, v := i, v
		go func() {
			fmt.Println(i, v)
		}()
	}

	// Prints 0 a, 1 b, 2 c
	for i, v := range vals {
		go func(i int, v string) {
			fmt.Println(i, v)
		}(i, v)
	}
}
