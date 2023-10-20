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
		return err
	}

	e, ok := v.Entries[cmd.ID]
	if !ok {
		return fmt.Errorf("no password for %q", cmd.ID)
	}

	fmt.Println(e.Password)
	return nil
}
