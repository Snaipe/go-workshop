package main

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"example.com/jaas/jq"
)

type Service struct {
	Timeout  time.Duration
	Cache    Cache[RunParams, []byte]
	CacheTTL time.Duration

	pool programPool
}

type programPool chan struct{}

func (p programPool) Acquire() {
	p <- struct{}{}
}

func (p programPool) Release() {
	<-p
}

func NewService(poolSize int) *Service {
	return &Service{
		pool: make(programPool, poolSize),
	}
}

type RunParams struct {
	Filter  string `form:"filter"`
	Data    string `form:"data"`
	Options string `form:"options"`
	Args    string `form:"args"`
}

func (svc *Service) Post(c echo.Context) error {
	var params RunParams
	c.Bind(&params)

	ctx, stop := context.WithTimeout(c.Request().Context(), svc.Timeout)
	defer stop()

	slog.Info("jq request",
		"filter", params.Filter,
		"data", params.Data,
		"options", params.Options,
		"args", params.Args,
	)

	cached, ok := svc.Cache.Get(params)
	if ok {
		return c.Blob(http.StatusOK, "application/json", cached)
	}

	filter := jq.NewFilter(params.Filter)

	opts := []jq.Option{
		jq.Context(ctx),
	}

	if params.Args != "" {
		for _, arg := range strings.Split(params.Args, ",") {
			name, value, found := strings.Cut(arg, "=")
			if !found {
				return echo.NewHTTPError(http.StatusBadRequest, "arg should be in the form key=value")
			}
			opts = append(opts, jq.Arg(name, value))
		}
	}

	var result bytes.Buffer

	svc.pool.Acquire()
	defer svc.pool.Release()

	err := filter.Run(strings.NewReader(params.Data), &result, opts...)

	var (
		syntaxErr *jq.SyntaxError
		parseErr  *jq.ParseError
	)
	switch {
	case errors.As(err, &syntaxErr):
		return echo.NewHTTPError(http.StatusBadRequest, err)
	case errors.As(err, &parseErr):
		return echo.NewHTTPError(http.StatusBadRequest, err)
	case errors.Is(err, context.DeadlineExceeded):
		return echo.NewHTTPError(http.StatusInternalServerError, "jq program timed out")
	case err != nil:
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	svc.Cache.Set(params, result.Bytes(), svc.CacheTTL)
	return c.Blob(http.StatusOK, "application/json", result.Bytes())
}
