package main

import (
	"fmt"
	"io"
	"crypto/sha256"
)

type DevZero struct{}

func (d DevZero) Read(buf []byte) (int, error) {
	clear(buf)
	return len(buf), nil
}

func main() {
	checksum := sha256.New()
	var zero DevZero

	io.Copy(checksum, &io.LimitedReader{R: zero, N: 128})

	fmt.Printf("%x\n", string(checksum.Sum(nil)))
}
