package main

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"

	"example.com/paman/vault"
)

type SetCmd struct {
	ID string `arg`
}

func updatePassword(id, password string, store vault.Store) error {
	var v vault.Vault
	err := store.Load(&v)
	switch {
	case errors.Is(err, os.ErrNotExist):
		// not fatal; create the vault
	case err != nil:
		return fmt.Errorf("loading vault: %w", err)
	}

	if v.Entries == nil {
		v.Entries = make(map[string]vault.Entry)
	}
	entry := v.Entries[id]
	entry.Password = password
	v.Entries[id] = entry

	if err := store.Store(&v); err != nil {
		return fmt.Errorf("storing vault: %w", err)
	}
	return nil
}

func (cmd *SetCmd) Run(store vault.Store) error {

	fmt.Fprint(os.Stderr, "please enter new password: ")
	password, err := term.ReadPassword(0)
	fmt.Fprint(os.Stderr, "\n")
	if err != nil {
		return fmt.Errorf("reading new password")
	}

	return updatePassword(cmd.ID, string(password), store)
}
