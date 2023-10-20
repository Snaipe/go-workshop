package jq

import (
	"errors"
	"io"
	"os/exec"
	"strings"
)

const (
	syntaxErrorCode = 3
	parseErrorCode  = 4
)

type Filter struct {
	filter string
}

func NewFilter(filter string) *Filter {
	return &Filter{filter: filter}
}

func (filter *Filter) Run(in io.Reader, out io.Writer, opts ...Option) error {
	var runopts runOptions
	for _, fn := range opts {
		fn(&runopts)
	}

	var stderr strings.Builder

	cmd := exec.CommandContext(runopts.Context, "jq", "-c", filter.filter)
	for k, v := range runopts.Args {
		cmd.Args = append(cmd.Args, "--argjson", k, v)
	}
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		var exiterr *exec.ExitError
		if errors.As(err, &exiterr) {
			switch exiterr.ExitCode() {
			case syntaxErrorCode:
				return &SyntaxError{Filter: filter.filter, Message: stderr.String()}
			case parseErrorCode:
				return &ParseError{Message: stderr.String()}
			}
			return errors.New(stderr.String())
		}
		return err
	}

	return nil
}

func (filter *Filter) String() string {
	return filter.filter
}
