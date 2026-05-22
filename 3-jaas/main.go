package main

import (
	"fmt"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const usage = `
jq as a service. Basic usage:
$ echo '{"a":"Hello","b":"World"}' | curl -F filter='.a+", "+.b+"!"' -F data=@- %s
"Hello, World!"
`

func main() {
	var cli struct {
		UsageAddress string `default:"localhost:1234"`
		BindAddress  string `arg optional default:":1234"`
		MaxProcs     int    `default:"1"`
		CacheSize    int    `default:"128"`
	}
	kong.Parse(&cli)

	e := echo.New()

	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	cache := NewInMemoryCache[string, []byte](cli.CacheSize)

	svc := NewService(cli.MaxProcs, cache)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf(usage, cli.UsageAddress))
	})
	e.POST("/", svc.Post)

	if err := e.Start(cli.BindAddress); err != nil {
		fmt.Println(err)
	}
}
