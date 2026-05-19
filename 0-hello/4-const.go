package main

func main() {
	const c = 123

	// Illegal
	//v := 1
	//const c2 = v + 2

	// OK
	const c2 = c + 2
}

const Constant = 1 << 256

type Animal int

const (
	Cat  Animal = iota // = 0
	Dog                // = 1
	Bird               // = 2
)

const (
	FlagA = 1 << iota // = 1
	FlagB             // = 2
	FlagC             // = 4
)
