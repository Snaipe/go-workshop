package main

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type CounterService struct {
	value int
	mux   sync.Mutex
}

func (svc *CounterService) GetCounter(c echo.Context) error {
	var result struct {
		Counter int `json:"counter"`
	}
	slog.Info("GetCounter")

	svc.mux.Lock()
	result.Counter = svc.value
	svc.mux.Unlock()
	return c.JSON(http.StatusOK, result)
}

func (svc *CounterService) PutCounter(c echo.Context) error {
	var params struct {
		Counter int `json:"counter"`
	}
	c.Bind(&params) // parse le corps de la requête HTTP dans `params`
	slog.Info("PutCounter", "params", params)

	svc.mux.Lock()
	svc.value = params.Counter
	svc.mux.Unlock()
	return c.JSON(http.StatusOK, params)
}

func (svc *CounterService) PatchCounter(c echo.Context) error {
	var params struct {
		Counter int `json:"counter"`
	}
	c.Bind(&params)

	svc.mux.Lock()
	svc.value += params.Counter
	params.Counter = svc.value
	svc.mux.Unlock()

	return c.JSON(http.StatusOK, params)
}
