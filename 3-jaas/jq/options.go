package jq

import (
	"context"
	"encoding/json"
	"strings"
)

type Option func(*runOptions)

type runOptions struct {
	Args    map[string]string
	Context context.Context
}

func Arg(name string, val string) Option {
	var encoded strings.Builder
	if err := json.NewEncoder(&encoded).Encode(val); err != nil {
		panic("programming error: cannot json-encode basic string")
	}
	return func(ro *runOptions) {
		if ro.Args == nil {
			ro.Args = map[string]string{}
		}
		ro.Args[name] = encoded.String()
	}
}

func ArgJSON(name, val string) Option {
	return func(ro *runOptions) {
		if ro.Args == nil {
			ro.Args = map[string]string{}
		}
		ro.Args[name] = val
	}
}

func Context(ctx context.Context) Option {
	return func(ro *runOptions) {
		ro.Context = ctx
	}
}
