package jq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestFilter(t *testing.T) {
	cases := []struct {
		Filter string
		Input  any
		Output string
		Args   map[string]string
	}{
		{
			Filter: ".title",
			Input: map[string]any{
				"userId":    1,
				"id":        1,
				"title":     "delectus aut autem",
				"completed": false,
			},
			Output: "\"delectus aut autem\"\n",
		},
		{
			Filter: ".",
			Input: map[string]any{
				"userId":    1,
				"id":        1,
				"title":     "delectus aut autem",
				"completed": false,
			},
			Output: `{"completed":false,"id":1,"title":"delectus aut autem","userId":1}` + "\n",
		},
		{
			Filter: ".[$key]",
			Input: map[string]any{
				"userId":    1,
				"id":        1,
				"title":     "delectus aut autem",
				"completed": false,
			},
			Output: "\"delectus aut autem\"\n",
			Args: map[string]string{
				"key": "title",
			},
		},
	}

	for i, tcase := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var in, out bytes.Buffer
			if err := json.NewEncoder(&in).Encode(tcase.Input); err != nil {
				t.Fatal(err)
			}
			opts := []Option{}
			for k, v := range tcase.Args {
				opts = append(opts, Arg(k, v))
			}
			if err := NewFilter(tcase.Filter).Run(&in, &out, opts...); err != nil {
				t.Fatal(err)
			}
			if actual := out.String(); tcase.Output != actual {
				t.Fatalf("expected %s, got %s", tcase.Output, actual)
			}
		})
	}
}

func TestFilterErrors(t *testing.T) {

	t.Run("SyntaxError", func(t *testing.T) {

		const in = `{"completed":false,"id":1,"title":"delectus aut autem","userId":1}`

		err := NewFilter("foo").Run(strings.NewReader(in), io.Discard)

		var synerr *SyntaxError
		if !errors.As(err, &synerr) {
			t.Fatal(err)
		}
	})

	t.Run("ParseError", func(t *testing.T) {
		const in = `bad`

		err := NewFilter(".").Run(strings.NewReader(in), io.Discard)

		var parseErr *ParseError
		if !errors.As(err, &parseErr) {
			t.Fatal(err)
		}
	})
}
