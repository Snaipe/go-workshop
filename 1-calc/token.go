package main

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenSymbol
	TokenNumber
	TokenLParen
	TokenRParen
)

type Token struct {
	Type TokenType
	Val  string
}

//go:generate go run golang.org/x/tools/cmd/stringer -type=TokenType -trimprefix Token

// "(+ 1 2
//         (+ foo 4))"

// "(", "+", "1", "2", "(", "+", "foo", "4", ")", ")"

// "(" -> LParen
// ")" -> RParen
// "+", "foo" -> Symbol
// "1", "2", "4" -> Number
// EOF
