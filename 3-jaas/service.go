package main

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"example.com/jaas/jq"
)

type Service struct{
	pool chan struct{}
}

func NewService(maxProcs int) *Service {
	return &Service{
		pool: make(chan struct{}, maxProcs),
	}
}

func (svc *Service) Post(c echo.Context) error {

	ctx, stop := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer stop()

	var params struct {
		Filter   string   `form:"filter"`
		Data     string   `form:"data"`
		Options  string   `form:"options"`
		Args     []string `form:"arg"`
		ArgJSONs []string `form:"argjson"`
	}
	c.Bind(&params)

	opts := []jq.Option{
		jq.Context(ctx),
	}
	for _, opt := range strings.Split(params.Options, " ") {
		switch opt {
		case "compact":
			opts = append(opts, jq.Compact())
		case "raw-output":
			opts = append(opts, jq.Raw())
		}
	}
	for _, arg := range params.Args {
		k, v, found := strings.Cut(arg, "=")
		if !found {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "arg does not have form key=value",
			})
		}
		opts = append(opts, jq.Arg(k, v))
	}
	for _, arg := range params.ArgJSONs {
		k, v, found := strings.Cut(arg, "=")
		if !found {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "arg does not have form key=value",
			})
		}
		opts = append(opts, jq.ArgJSON(k, v))
	}

	filter := jq.NewFilter(params.Filter)

	select {
	case svc.pool <- struct{}{}:
	case <-ctx.Done():
		return c.NoContent(http.StatusTooManyRequests)
	}
	defer func() {
		<-svc.pool
	}()

	slog.Info("POST", "params", params)
	in := strings.NewReader(params.Data)
	var out bytes.Buffer
	if err := filter.Run(in, &out, opts...); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err})
	}

	return c.Blob(http.StatusOK, "application/json", out.Bytes())
}
