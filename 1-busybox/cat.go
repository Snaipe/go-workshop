package main

import (
	"bufio"
	"fmt"
	"os"
)

type CatCmd struct {
	Number bool `short:"n" help:"display line numbers"`
}

func (cmd *CatCmd) Run() error {
	scanner := bufio.NewScanner(os.Stdin)
	lineno := 1
	for scanner.Scan() {
		if cmd.Number {
			fmt.Printf("%d\t%s\n", lineno, scanner.Text())
		} else {
			fmt.Println(scanner.Text())
		}
		lineno++
	}
	return scanner.Err()
}
