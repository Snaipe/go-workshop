package main

import (
	"fmt"
	"strings"
)

type EchoCmd struct {
	NoNewline bool     `short:"n" name:""`
	Args      []string `arg:"" passthrough:""`
}

func (cmd *EchoCmd) Run() error {
	out := strings.Join(cmd.Args, " ")
	if !cmd.NoNewline {
		out += "\n"
	}
	_, err := fmt.Print(out)
	return err
}
