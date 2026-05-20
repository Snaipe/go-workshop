package main

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

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
