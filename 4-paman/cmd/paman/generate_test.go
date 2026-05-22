package main

import (
	"testing"
)

func TestParseRange(t *testing.T) {

	type TestCase struct {
		Input string
		Set   string
	}

	tcases := []TestCase{
		{
			Input: "a",
			Set:   "a",
		},
		{
			Input: "a-",
			Set:   "-a",
		},
		{
			Input: "-a",
			Set:   "-a",
		},
		{
			Input: "a-f",
			Set:   "abcdef",
		},
		{
			Input: "a-fA-F",
			Set:   "ABCDEFabcdef",
		},
		{
			Input: "a-fA-F ",
			Set:   " ABCDEFabcdef",
		},
	}

	for _, tcase := range tcases {
		t.Run(tcase.Input, func(t *testing.T) {
			set := parseRange(tcase.Input)
			if set != tcase.Set {
				t.Fatalf("expected %q, got %q", tcase.Set, set)
			}
		})
	}
}
