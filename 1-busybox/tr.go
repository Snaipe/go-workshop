package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type TrCmd struct {
	Squeeze bool   `short:"s"`
	Delete  bool   `short:"d"`
	Set1    string `arg`
	Set2    string `arg optional`
}

func (cmd *TrCmd) Run() error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case cmd.Squeeze:
			for {
				idx := strings.Index(line, cmd.Set1)
				if idx == -1 {
					break
				}
				fmt.Print(line[:idx+1])
				line = strings.TrimLeft(line[idx+1:], cmd.Set1)
			}
			fmt.Println(line)
		case cmd.Delete:
			fmt.Println(strings.ReplaceAll(line, cmd.Set1, ""))
		default:
			if cmd.Set2 == "" {
				return fmt.Errorf("Set2 must be non-empty")
			}
			fmt.Println(strings.ReplaceAll(line, cmd.Set1, cmd.Set2))
		}
	}
	return scanner.Err()
}
