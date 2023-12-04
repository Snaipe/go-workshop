package main

import (
	"github.com/alecthomas/kong"
)

func main() {
	var cli struct {
		Echo EchoCmd `cmd:""`
		Cp   CpCmd   `cmd:""`
		Cat  CatCmd  `cmd:""`
		Ls   LsCmd   `cmd:""`
		Du   DuCmd   `cmd:""`
		Nc   NcCmd   `cmd:""`
	}
	ctx := kong.Parse(&cli, kong.UsageOnError())
	ctx.FatalIfErrorf(ctx.Run())
}
