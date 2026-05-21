package jq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Options struct {
	Compact bool
	Raw     bool
	Args    map[string]string
	Context context.Context
}

type Option func(*Options) error

func Compact() Option {
	return func(opts *Options) error {
		opts.Compact = true
		return nil
	}
}

func Raw() Option {
	return func(opts *Options) error {
		opts.Raw = true
		return nil
	}
}

func Context(ctx context.Context) Option {
	return func(opts *Options) error {
		opts.Context = ctx
		return nil
	}
}

func Arg(key, value string) Option {
	return func(opts *Options) error {
		if opts.Args == nil {
			opts.Args = make(map[string]string)
		}
		txt, err := json.Marshal(value)
		if err != nil {
			return err
		}
		opts.Args[key] = string(txt)
		return nil
	}
}

func ArgJSON(key, value string) Option {
	return func(opts *Options) error {
		if opts.Args == nil {
			opts.Args = make(map[string]string)
		}
		opts.Args[key] = value
		return nil
	}
}

// Filter is a jqlang program (called a "filter") that is used to process
// JSON data.
type Filter struct {
	filter string
}

func NewFilter(filter string) *Filter {
	return &Filter{filter: filter}
}

// Run executes the filter on the contents of in, and writes the result
// to out.
func (f *Filter) Run(in io.Reader, out io.Writer, opts ...Option) error {

	var o Options
	for _, opt := range opts {
		if err := opt(&o); err != nil {
			return err
		}
	}

	args := []string{f.filter}
	if o.Compact {
		args = append(args, "-c")
	}
	if o.Raw {
		args = append(args, "-r")
	}
	for k, v := range o.Args {
		args = append(args, "--argjson", k, v)
	}

	var stderr strings.Builder

	cmd := exec.CommandContext(o.Context, "jq", args...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if err := o.Context.Err(); err != nil {
			return err
		}
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
