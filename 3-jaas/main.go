package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alecthomas/kong"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const usage = `
     _             ____
    | | __ _  __ _/ ___|
 _  | |/ _\ |/ _\ \___ \
| |_| | (_| | (_| |___) |
 \___/ \__,_|\__,_|____/

jq as a service. Basic usage:
$ echo '{"a":"Hello","b":"World"}' | curl -F filter='.a+", "+.b+"!"' -F data=@- %s
"Hello, World!"
`

func main() {
	var cli struct {
		UsageAddress string        `default:"localhost:1234"`
		Timeout      time.Duration `default:"1m"`
		PoolSize     int           `default:"1024"`
		CacheTTL     time.Duration `default:"1h"`

		BindAddress string `arg:"" optional:"" default:":1234"`
	}
	kong.Parse(&cli)

	e := echo.New()

	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	svc := NewService(cli.PoolSize)
	svc.Timeout = cli.Timeout
	svc.Cache = NewInMemoryCache[RunParams, []byte]()
	svc.CacheTTL = cli.CacheTTL

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf(usage, cli.UsageAddress))
	})
	e.POST("/", svc.Post)

	if err := e.Start(cli.BindAddress); err != nil {
		fmt.Println(err)
	}
}
