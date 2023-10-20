package vault

import (
	"crypto/cipher"
	"errors"
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
	switch {
	case errors.Is(err, os.ErrNotExist):
		*vault = Vault{
			Entries: map[string]Entry{},
		}
		return nil
	case err != nil:
		return err
	}
	defer f.Close()

	return vault.Unmarshal(f, store.Block)
}

func (store *FileStore) Store(vault *Vault) error {
	f, err := os.Create(store.Path + ".tmp")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := vault.Marshal(f, store.Block); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return os.Rename(store.Path+".tmp", store.Path)
}
