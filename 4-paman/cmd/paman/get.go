package main

import (
	"fmt"

	"example.com/paman/vault"
)

type GetCmd struct {
	ID string `arg`
}

func (cmd *GetCmd) Run(store vault.Store) error {

	var v vault.Vault
	if err := store.Load(&v); err != nil {
		return fmt.Errorf("loading vault: %w", err)
	}

	entry, found := v.Entries[cmd.ID]
	if !found {
		return fmt.Errorf("Missing vault entry for %q", cmd.ID)
	}

	fmt.Printf("password: %q\n", entry.Password)
	return nil
}
