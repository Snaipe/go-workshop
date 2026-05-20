package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	// Créée une instance d'Echo
	e := echo.New()

	// Associe une fonction qui retourne un HTTP 200 OK avec en corps
	// "Hello, World!\n" à la racine du service
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	var svc CounterService
	e.GET("/counter", svc.GetCounter)
	e.PUT("/counter", svc.PutCounter)
	e.PATCH("/counter", svc.PatchCounter)

	// Démarre le serveur et bloque jusqu'à ce qu'il finisse
	if err := e.Start(":1323"); err != nil {
		fmt.Println(err)
	}
}
