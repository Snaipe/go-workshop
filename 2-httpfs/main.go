package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/labstack/echo/v4"
)

func main() {
	var cli struct {
		BindAddress string `arg`
		RootDir     string `arg`
	}
	kong.Parse(&cli)

	e := echo.New()

	svc := Service{
		RootDir: cli.RootDir,
	}
	e.GET("/*", svc.GetPath)
	e.PUT("/*", svc.PutPath)

	if err := e.Start(cli.BindAddress); err != nil {
		fmt.Println(err)
	}
}
