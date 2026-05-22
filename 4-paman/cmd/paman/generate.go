package main

import (
	"math/rand/v2"

	"example.com/paman/vault"
)

type GenerateCmd struct {
	Chars  string `default:"a-zA-Z0-9!-/:-@[-"`
	ID     string `arg`
	Length int    `arg`
}

func parseRange(format string) string {

	var prec byte

	next := func() (byte, bool) {
		if prec != 0 {
			c := prec
			prec = 0
			return c, true
		}
		if len(format) == 0 {
			return 0, false
		}
		c := format[0]
		format = format[1:]
		return c, true
	}

	back := func(b byte) {
		prec = b
	}

	var set [128]byte

	for {
		start, ok := next()
		if !ok {
			break
		}

		maybeDash, ok := next()
		if ok && maybeDash == '-' {
			end, ok := next()
			if !ok {
				set[start] = start
				set[maybeDash] = maybeDash
				break
			} else {
				for i := start; i <= end; i++ {
					set[i] = i
				}
			}
			continue
		}
		if ok {
			back(maybeDash)
		}

		set[start] = start
	}

	var reduced []byte
	for _, b := range set {
		if b != 0 {
			reduced = append(reduced, b)
		}
	}
	return string(reduced)
}

func (cmd *GenerateCmd) Run(store vault.Store) error {

	set := parseRange(cmd.Chars)

	password := make([]byte, cmd.Length)
	for i := range cmd.Length {
		password[i] = set[rand.N(len(set))]
	}

	return updatePassword(cmd.ID, string(password), store)
}
