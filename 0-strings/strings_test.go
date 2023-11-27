package strings

import (
	"slices"
	"testing"
)

func TestIndexByte(t *testing.T) {

	type TestCase struct {
		In       string
		Char     byte
		Expected int
	}

	tcases := []TestCase{
		{
			In:       "",
			Char:     0,
			Expected: -1,
		},
		{
			In:       "abcd",
			Char:     'b',
			Expected: 1,
		},
		{
			In:       "abcd",
			Char:     'e',
			Expected: -1,
		},
		{
			In:       "Hello, 世界",
			Char:     184,
			Expected: 8,
		},
	}

	for _, tcase := range tcases {
		t.Run(tcase.In, func(t *testing.T) {
			if actual := IndexByte(tcase.In, tcase.Char); tcase.Expected != actual {
				t.Fatalf("expected %v, got %v", tcase.Expected, actual)
			}
		})
	}
}

func TestIndex(t *testing.T) {

	type TestCase struct {
		In       string
		Substr   string
		Expected int
	}

	tcases := []TestCase{
		{
			In:       "",
			Substr:   "abc",
			Expected: -1,
		},
		{
			In:       "abcd",
			Substr:   "",
			Expected: 0,
		},
		{
			In:       "abcd",
			Substr:   "e",
			Expected: -1,
		},
		{
			In:       "abcd",
			Substr:   "ac",
			Expected: -1,
		},
		{
			In:       "Hello, 世界",
			Substr:   "界",
			Expected: 10,
		},
	}

	for _, tcase := range tcases {
		t.Run(tcase.In, func(t *testing.T) {
			if actual := Index(tcase.In, tcase.Substr); tcase.Expected != actual {
				t.Fatalf("expected %v, got %v", tcase.Expected, actual)
			}
		})
	}
}

func TestCut(t *testing.T) {

	type TestCase struct {
		In     string
		Sep    string
		Before string
		After  string
		Found  bool
	}

	tcases := []TestCase{
		{
			In:     "",
			Sep:    "",
			Before: "",
			After:  "",
			Found:  true,
		},
		{
			In:     "abc",
			Sep:    "",
			Before: "",
			After:  "abc",
			Found:  true,
		},
		{
			In:     "abc",
			Sep:    "b",
			Before: "a",
			After:  "c",
			Found:  true,
		},
		{
			In:     "abc",
			Sep:    "e",
			Before: "abc",
			After:  "",
			Found:  false,
		},
	}

	for _, tcase := range tcases {
		t.Run(tcase.In, func(t *testing.T) {
			before, after, found := Cut(tcase.In, tcase.Sep)
			if before != tcase.Before {
				t.Fatalf("expected before %v, got %v", tcase.Before, before)
			}
			if after != tcase.After {
				t.Fatalf("expected after %v, got %v", tcase.After, after)
			}
			if found != tcase.Found {
				t.Fatalf("expected found %v, got %v", tcase.Found, found)
			}
		})
	}
}

func TestSplit(t *testing.T) {

	type TestCase struct {
		In       string
		Sep      string
		Expected []string
	}

	tcases := []TestCase{
		{
			In:       "",
			Sep:      "",
			Expected: []string{},
		},
		{
			In:       "abc",
			Sep:      "",
			Expected: []string{"a", "b", "c"},
		},
		{
			In:       "aa:b:c",
			Sep:      ":",
			Expected: []string{"aa", "b", "c"},
		},
		{
			In:       "aa:b:c",
			Sep:      ",",
			Expected: []string{"aa:b:c"},
		},
	}

	for _, tcase := range tcases {
		t.Run(tcase.In, func(t *testing.T) {
			actual := Split(tcase.In, tcase.Sep)
			if !slices.Equal(actual, tcase.Expected) {
				t.Fatalf("expected %#v, got %#v", tcase.Expected, actual)
			}
		})
	}
}
