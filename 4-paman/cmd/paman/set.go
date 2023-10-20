package main

import (
	"fmt"
	"os"

	"golang.org/x/term"

	"example.com/paman/vault"
)

type SetCmd struct {
	ID string `arg`
}

func setPassword(store vault.Store, id, password string) error {
	var v vault.Vault
	if err := store.Load(&v); err != nil {
		return err
	}

	entry := v.Entries[id]
	entry.Password = password
	v.Entries[id] = entry

	return store.Store(&v)
}

func (cmd *SetCmd) Run(store vault.Store) error {
	fmt.Fprint(os.Stderr, "Please enter the password: ")
	password, err := term.ReadPassword(0)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr)

	return setPassword(store, cmd.ID, string(password))
}
