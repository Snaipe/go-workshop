package main

import (
	"crypto/rand"
	"fmt"
	"slices"

	"example.com/paman/vault"
)

type GenerateCmd struct {
	Chars  string `default:"a-zA-Z0-9!-/:-@[-"`
	ID     string `arg`
	Length int    `arg`
}

func parseCharSet(s string) []rune {
	idx := 0

	next := func() rune {
		if idx >= len(s) {
			return 0
		}
		r := s[idx]
		idx++
		return rune(r)
	}

	back := func() {
		if idx > 0 {
			idx--
		}
	}

	peek := func() rune {
		r := next()
		if r != 0 {
			back()
		}
		return r
	}

	var allowedChars []rune
	for {
		r := next()
		if r == 0 {
			break
		}

		if r2 := peek(); r2 == '-' {
			next()
			r3 := peek()
			if r3 != 0 {
				next()
				from := r
				to := r3
				for r := from; r <= to; r++ {
					allowedChars = append(allowedChars, r)
				}
				continue
			}
			allowedChars = append(allowedChars, r2)
		}
		allowedChars = append(allowedChars, r)
	}

	slices.Sort(allowedChars)
	return slices.Compact(allowedChars)
}

func (cmd *GenerateCmd) Run(store vault.Store) error {
	set := parseCharSet(cmd.Chars)

	password := make([]byte, cmd.Length)
	if _, err := rand.Read(password); err != nil {
		return err
	}
	for i, b := range password {
		password[i] = byte(set[int(b)%len(set)])
	}

	strPassword := string(password)
	fmt.Println(strPassword)
	return setPassword(store, cmd.ID, strPassword)
}
