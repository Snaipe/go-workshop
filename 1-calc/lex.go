package main

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	in        string
	start     int
	index     int
	lastWidth int
	tokens    []Token
}

func MakeLexer(in string) Lexer {
	return Lexer{
		in: in,
	}
}

func (l *Lexer) next() rune {
	if l.index >= len(l.in) {
		l.lastWidth = 0
		return 0
	}
	r, w := utf8.DecodeRuneInString(l.in[l.index:])
	l.lastWidth = w
	l.index += w
	return r
}

func (l *Lexer) back() {
	l.index -= l.lastWidth
	l.lastWidth = 0
}

func (l *Lexer) peek() rune {
	r := l.next()
	if r != 0 {
		l.back()
	}
	return r
}

func (l *Lexer) emit(typ TokenType) {
	val := l.in[l.start:l.index]
	l.start = l.index
	l.lastWidth = 0
	l.tokens = append(l.tokens, Token{Type: typ, Val: val})
}

func (l *Lexer) discard() {
	l.start = l.index
}

type stateFunc func() stateFunc

func (l *Lexer) Lex() []Token {
	state := l.lexMain
	for state != nil {
		state = state()
	}
	return l.tokens
}

func (l *Lexer) lexMain() stateFunc {
	for {
		n := l.next()
		switch {
		case n == 0: // EOF
			return nil // On arrête la machine à état
		case n == '(':
			l.emit(TokenLParen)
		case n == ')':
			l.emit(TokenRParen)
		case n == '-':
			if strings.IndexRune("0123456789", l.peek()) == -1 {
				// n est un tiret qui n'est pas suivi d'un
				// chiffre; c'est le début d'un symbole.
				return l.lexSymbol
			}
			fallthrough
		case n >= '0' && n <= '9':
			return l.lexNumber
		case unicode.IsSpace(n):
			l.discard() // On exclue les espaces de nos jetons
		default:
			return l.lexSymbol
		}
	}
}

func (l *Lexer) lexSymbol() stateFunc {
	for {
		r := l.next()
		if r == 0 || unicode.IsSpace(r) || strings.IndexRune("()0123456789", r) != -1 {
			break
		}
	}
	l.back()
	l.emit(TokenSymbol)
	return l.lexMain
}

func (l *Lexer) lexNumber() stateFunc {
	for {
		r := l.next()
		if r == 0 || strings.IndexRune("0123456789", r) == -1 {
			break
		}
	}
	l.back()
	l.emit(TokenNumber)
	return l.lexMain
}

// "(+ 1 2  (+ foo 4))"
// ["f", "o", "o", ]" ", "4", ")", ")"
//
// lex.next()
// lex.emit(TokenLParen)
// lex.next()
// lex.peek()
// lex.emit(TokenSymbol)
// lex.next()
// lex.discard()
// lex.next()
// lex.peek()
// lex.emit(TokenNumber)
// ...
// lex.next()
// lex.next()
// lex.next()
// lex.emit(TokenSymbol) // Token{Type: TokenSymbol, Val: "foo"}

// "(", "+", "1", "2", "(", "+", "foo", "4", ")", ")"
