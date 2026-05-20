package main

import (
	"github.com/alecthomas/kong"
)

func main() {
	var cli struct {
		Echo EchoCmd `cmd:""`
		Cat  CatCmd  `cmd:""`
		Cp   CpCmd   `cmd:""`
		Ls   LsCmd   `cmd:""`
		Du   DuCmd   `cmd:""`
		Nc   NcCmd   `cmd:""`
	}
	ctx := kong.Parse(&cli)
	ctx.FatalIfErrorf(ctx.Run())
}
