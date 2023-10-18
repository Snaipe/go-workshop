package main

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type CounterService struct {
	value int
	mux   sync.RWMutex
}

func (svc *CounterService) GetCounter(c echo.Context) error {
	var result struct {
		Counter int `json:"counter"`
	}
	svc.mux.RLock()
	result.Counter = svc.value
	svc.mux.RUnlock()
	return c.JSON(http.StatusOK, result)
}

func (svc *CounterService) PutCounter(c echo.Context) error {
	var result struct {
		Counter int `json:"counter"`
	}
	c.Bind(&result)
	svc.mux.Lock()
	svc.value = result.Counter
	svc.mux.Unlock()
	return c.JSON(http.StatusOK, result)
}

func (svc *CounterService) PatchCounter(c echo.Context) error {
	var result struct {
		Counter int `json:"counter"`
	}
	c.Bind(&result)
	svc.mux.Lock()
	svc.value += result.Counter
	result.Counter = svc.value
	svc.mux.Unlock()
	return c.JSON(http.StatusOK, result)
}
