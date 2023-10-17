package main

import (
	"fmt"
	"strings"
)

type EchoCmd struct {
	Args      []string `arg optional`
	NoNewline bool     `short:"n" help:"disable newlines"`
}

func (cmd *EchoCmd) Run() error {
	fmt.Print(strings.Join(cmd.Args, " "))
	if !cmd.NoNewline {
		fmt.Println()
	}
	return nil
}
