package strings

import (
	"testing"
)

func TestIndexByte(t *testing.T) {

	type TestCase struct {
		Input string
		Byte  byte
		Index int
	}

	cases := []TestCase{
		{
			Input: "",
			Byte: 0,
			Index: -1,
		},
		{
			Input: "bonjour",
			Byte: 'j',
			Index: 3,
		},
		{
			Input: "bonjour",
			Byte: 'x',
			Index: -1,
		},
		{
			Input: "Hello, 世界",
			Byte: 184,
			Index: 8,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			got := IndexByte(tc.Input, tc.Byte)
			if got != tc.Index {
				t.Fatalf("wanted %d, got %d", tc.Index, got)
			}
		})
	}
}

func TestIndex(t *testing.T) {

	type TestCase struct {
		Input  string
		Substr string
		Index  int
	}

	cases := []TestCase{
		{
			Input:  "Bonjour",
			Substr: "",
			Index:  0,
		},
		{
			Input:  "Bonjour",
			Substr: "jour",
			Index:  3,
		},
		{
			Input:  "Bonjour",
			Substr: "joux",
			Index:  -1,
		},
		{
			Input:  "Hello, 世界",
			Substr: "世",
			Index:  7,
		},
		{
			Input:  "Bonjour",
			Substr: "ou",
			Index:  4,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			got := Index(tc.Input, tc.Substr)
			if got != tc.Index {
				t.Fatalf("wanted %d, got %d", tc.Index, got)
			}
		})
	}
}

func TestCut(t *testing.T) {

	type TestCase struct {
		Input  string
		Sep    string
		Before string
		After  string
		Found  bool
	}

	cases := []TestCase{
		{
			Input:  "Bonjour",
			Sep:    "",
			Before: "",
			After:  "Bonjour",
			Found:  true,
		},
		{
			Input:  "Bonjour",
			Sep:    "ou",
			Before: "Bonj",
			After:  "r",
			Found:  true,
		},
		{
			Input:  "a::b:",
			Sep:    ":",
			Before: "a",
			After:  ":b:",
			Found:  true,
		},
		{
			Input:  "a::b:",
			Sep:    "x",
			Before: "a::b:",
			After:  "",
			Found:  false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			before, after, found := Cut(tc.Input, tc.Sep)
			if before != tc.Before {
				t.Fatalf("before: wanted %v, got %v", tc.Before, before)
			}
			if after != tc.After {
				t.Fatalf("after: wanted %v, got %v", tc.After, after)
			}
			if found != tc.Found {
				t.Fatalf("found: wanted %v, got %v", tc.Found, found)
			}
		})
	}
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestSplit(t *testing.T) {

	type TestCase struct {
		Input string
		Sep   string
		Want  []string
	}

	cases := []TestCase{
		{
			Input: "a::b:",
			Sep:   ":",
			Want:  []string{"a", "", "b", ""},
		},
		{
			Input: "a::b:",
			Sep:   "x",
			Want:  []string{"a::b:"},
		},
		{
			Input: "a::b:",
			Sep:   "",
			Want:  []string{"a", ":", ":", "b", ":"},
		},
		{
			Input: "",
			Sep:   "",
			Want:  []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			got := Split(tc.Input, tc.Sep)
			if !slicesEqual(got, tc.Want) {
				t.Fatalf("wanted %#v, got %#v", tc.Want, got)
			}
		})
	}
}

func TestIndexRune(t *testing.T) {

	type TestCase struct {
		Input string
		Rune  rune
		Index int
	}

	cases := []TestCase{
		{
			Input: "",
			Rune: 0,
			Index: -1,
		},
		{
			Input: "bonjour",
			Rune: 'j',
			Index: 3,
		},
		{
			Input: "bonjour",
			Rune: 'x',
			Index: -1,
		},
		{
			Input: "Hello, 世界",
			Rune: '界',
			Index: 10,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			got := IndexRune(tc.Input, tc.Rune)
			if got != tc.Index {
				t.Fatalf("wanted %d, got %d", tc.Index, got)
			}
		})
	}
}

func TestIndexAny(t *testing.T) {

	type TestCase struct {
		Input string
		Chars string
		Index int
	}

	cases := []TestCase{
		{
			Input: "abc123",
			Chars: "1234567890",
			Index: 3,
		},
		{
			Input: "bonjour",
			Chars: "xr",
			Index: 6,
		},
		{
			Input: "bonjour",
			Chars: "",
			Index: -1,
		},
		{
			Input: "Hello, 世界, world",
			Chars: "界世",
			Index: 7,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			got := IndexAny(tc.Input, tc.Chars)
			if got != tc.Index {
				t.Fatalf("wanted %d, got %d", tc.Index, got)
			}
		})
	}
}
