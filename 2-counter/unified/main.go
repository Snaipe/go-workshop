package main

import (
	"github.com/alecthomas/kong"
)

func main() {
	var cli struct {
		Serve ServeCmd `cmd:""`
		Get   GetCmd   `cmd:""`
		Set   SetCmd   `cmd:""`
		Add   AddCmd   `cmd:""`
	}

	ctx := kong.Parse(&cli)
	ctx.FatalIfErrorf(ctx.Run())
}
