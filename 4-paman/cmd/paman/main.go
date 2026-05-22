package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

func main() {
	var cli struct {
		Get      GetCmd      `cmd help:"retrieve a password"`
		Set      SetCmd      `cmd help:"assign a password"`
		Generate GenerateCmd `cmd help:"generate a password"`
	}

	ctx := kong.Parse(&cli)

	if err := ctx.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		os.Exit(1)
	}
}
