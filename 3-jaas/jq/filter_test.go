package jq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestFilter(t *testing.T) {
	cases := []struct {
		Filter  string
		Input   any
		Output  string
		Options []Option
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
			Filter: "to_entries",
			Input: map[string]any{
				"userId":    1,
				"id":        1,
				"title":     "delectus aut autem",
				"completed": false,
			},
			Output:  `[{"key":"completed","value":false},{"key":"id","value":1},{"key":"title","value":"delectus aut autem"},{"key":"userId","value":1}]` + "\n",
			Options: []Option{Compact()},
		},
	}

	for i, tcase := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var in, out bytes.Buffer
			if err := json.NewEncoder(&in).Encode(tcase.Input); err != nil {
				t.Fatal(err)
			}
			if err := NewFilter(tcase.Filter).Run(&in, &out, tcase.Options...); err != nil {
				t.Fatal(err)
			}
			if actual := out.String(); tcase.Output != actual {
				t.Fatalf("expected %s, got %s", tcase.Output, actual)
			}
		})
	}
}

func TestFilterErrors(t *testing.T) {
	cases := []struct {
		Filter string
		Input  string
		Error  string
	}{
		{
			Filter: ".title %% notexist",
			Input:  "",
			Error:  "syntax error on filter \".title %% notexist\": jq: error: syntax error, unexpected '%' at <top-level>, line 1, column 9:\n    .title %% notexist\n            ^\njq: 1 compile error\n",
		},
		{
			Filter: ".title",
			Input:  "{",
			Error:  "json parse error: jq: parse error: Unfinished JSON term at EOF at line 1, column 1\n",
		},
	}

	for i, tcase := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var out bytes.Buffer
			in := strings.NewReader(tcase.Input)
			err := NewFilter(tcase.Filter).Run(in, &out)

			if err == nil {
				t.Fatalf("expected error %q, got nil", tcase.Error)
			}
			if got := err.Error(); got != tcase.Error {
				t.Fatalf("expected error %q, got %q instead", tcase.Error, got)
			}
		})
	}
}
