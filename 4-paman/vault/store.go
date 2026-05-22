package vault

import (
	"crypto/cipher"
	"os"
)

type Store interface {
	Load(*Vault) error
	Store(*Vault) error
}

type FileStore struct {
	Block cipher.Block
	Path  string
}

func (store *FileStore) Load(vault *Vault) error {
	f, err := os.Open(store.Path)
	if err != nil {
		return err
	}

	return vault.Unmarshal(f, store.Block)
}

func (store *FileStore) Store(vault *Vault) error {
	f, err := os.Create(store.Path + ".new")
	if err != nil {
		return err
	}

	if err := vault.Marshal(f, store.Block); err != nil {
		return err
	}

	if err := os.Rename(f.Name(), store.Path); err != nil {
		return err
	}

	return nil
}
