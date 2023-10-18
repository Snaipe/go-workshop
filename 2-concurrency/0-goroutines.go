package main

import (
	"fmt"
	"time"
)

func main() {
	// La fonction est lancée de manière concurrente
	go func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("first!")
	}()

	time.Sleep(1 * time.Second)
	fmt.Println("second!")
}
