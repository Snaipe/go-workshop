package main

import (
	"fmt"
	"slices"
	"testing"
)

func TestParseCharset(t *testing.T) {
	cases := []struct {
		Input    string
		Expected string
	}{
		{
			Input:    "a-z",
			Expected: "abcdefghijklmnopqrstuvwxyz",
		},
		{
			Input:    "a-z-",
			Expected: "-abcdefghijklmnopqrstuvwxyz",
		},
		{
			Input:    "-a-z",
			Expected: "-abcdefghijklmnopqrstuvwxyz",
		},
		{
			Input:    "a-zA-Z0-9!-/:-@[-",
			Expected: "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[abcdefghijklmnopqrstuvwxyz",
		},
	}

	for i, tcase := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			expected := []rune(tcase.Expected)
			actual := parseCharSet(tcase.Input)

			if !slices.Equal(expected, actual) {
				t.Fatalf("expected %q, got %q", string(expected), string(actual))
			}
		})
	}
}
