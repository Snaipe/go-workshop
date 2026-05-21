package jq

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Options struct {
	Compact bool
}

type Option func(*Options)

func Compact() Option {
	return func(opts *Options) {
		opts.Compact = true
	}
}

// Filter is a jqlang program (called a "filter") that is used to process
// JSON data.
type Filter struct {
	filter string
	opts   Options
}

func NewFilter(filter string, opts ...Option) *Filter {
	f := &Filter{filter: filter}
	for _, opt := range opts {
		opt(&f.opts)
	}
	return f
}

// Run executes the filter on the contents of in, and writes the result
// to out.
func (f *Filter) Run(in io.Reader, out io.Writer) error {

	args := []string{f.filter}
	if f.opts.Compact {
		args = append(args, "-c")
	}

	var stderr strings.Builder

	cmd := exec.Command("jq", args...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		var exiterr *exec.ExitError
		if errors.As(err, &exiterr) {
			switch exiterr.ExitCode() {
			case 3:
				return &SyntaxError{
					Filter:  f.filter,
					Message: stderr.String(),
				}
			case 5:
				return &ParseError{
					Message: stderr.String(),
				}
			}
			return fmt.Errorf("system error: %v", stderr.String())
		}
	}

	return err
}

func (f *Filter) String() string {
	return f.filter
}
