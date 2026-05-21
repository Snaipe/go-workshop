package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type Service struct{}

func (svc *Service) Post(c echo.Context) error {
	var params struct {
		Filter  string `form:"filter"`
		Data    string `form:"data"`
		Options string `form:"options"`
		Args    string `form:"args"`
	}
	c.Bind(&params)
	return fmt.Errorf("unimplemented")
}
