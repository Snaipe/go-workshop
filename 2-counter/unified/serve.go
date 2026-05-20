package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type ServeCmd struct {
	BindAddress string `arg:"" default:":1323"`
}

func (c *ServeCmd) Run() error {

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
	if err := e.Start(c.BindAddress); err != nil {
		fmt.Println(err)
	}
	return nil
}

type CounterService struct {
	value int
	mu    sync.Mutex
}

func (svc *CounterService) GetCounter(c echo.Context) error {
	var result struct {
		Counter int `json:"counter"`
	}

	svc.mu.Lock()
	result.Counter = svc.value
	svc.mu.Unlock()

	return c.JSON(http.StatusOK, result)
}

func (svc *CounterService) PutCounter(c echo.Context) error {
	var params struct {
		Counter int `json:"counter"`
	}
	c.Bind(&params) // parse le corps de la requête HTTP dans `params`

	svc.mu.Lock()
	svc.value = params.Counter
	svc.mu.Unlock()

	return c.JSON(http.StatusOK, params)
}

func (svc *CounterService) PatchCounter(c echo.Context) error {
	var params struct {
		Counter int `json:"counter"`
	}
	c.Bind(&params) // parse le corps de la requête HTTP dans `params`

	svc.mu.Lock()
	svc.value += params.Counter
	svc.mu.Unlock()

	return c.NoContent(http.StatusOK)
}
