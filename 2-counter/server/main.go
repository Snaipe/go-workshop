package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	//sigchan := make(chan os.Signal)
	//signal.Notify(sigchan, syscall.SIGINT)
	//go func() {
	//	
	//}()

	var svc CounterService
	e.GET("/counter", svc.GetCounter)
	e.PUT("/counter", svc.PutCounter)
	e.PATCH("/counter", svc.PatchCounter)

	if err := e.Start(":1323"); err != nil {
		fmt.Println(err)
	}
}
