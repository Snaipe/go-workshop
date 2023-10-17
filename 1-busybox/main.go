package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

type CLI struct {
	Echo EchoCmd `cmd`
	Cat  CatCmd  `cmd`
	Tr   TrCmd   `cmd`
}

func main() {

	var cli CLI

	ctx := kong.Parse(&cli)
	if err := ctx.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
