package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

type ServeCmd struct {
	BindAddress string `default:":1323" help:"the address to bind to"`
}

func (cmd *ServeCmd) Run() error {

	// Créée une instance d'Echo
	e := echo.New()

	var svc CounterService
	e.GET("/counter", svc.GetCounter)
	e.PUT("/counter", svc.PutCounter)
	e.PATCH("/counter", svc.PatchCounter)

	// Démarre le serveur et bloque jusqu'à ce qu'il finisse
	if err := e.Start(cmd.BindAddress); err != nil {
		fmt.Println(err)
	}
	return nil
}
