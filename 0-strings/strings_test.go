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
	panic("unimplemented")
}

func TestCut(t *testing.T) {
	panic("unimplemented")
}

func TestSplit(t *testing.T) {
	panic("unimplemented")
}
