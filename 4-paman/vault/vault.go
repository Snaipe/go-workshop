package vault

import (
	"encoding/json"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type Entry struct {
	Password string
	// Potentiellement d'autres champs plus tard
}

type Vault struct {
	Entries map[string]Entry
}

func (v *Vault) Marshal(out io.Writer, block cipher.Block) error {
	plaintext, err := json.Marshal(v.Entries)
	if err != nil {
		return err
	}

	blocksize := block.BlockSize()

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, blocksize+len(plaintext))
	iv := ciphertext[:blocksize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[blocksize:], plaintext)

	n, err := out.Write(ciphertext)
	if err != nil {
		return err
	}
	if n != len(ciphertext) {
		return io.ErrShortWrite
	}

	return nil
}

func (v *Vault) Unmarshal(in io.Reader, block cipher.Block) error {
	blocksize := block.BlockSize()
	iv := make([]byte, blocksize)

	if _, err := io.ReadFull(in, iv); err != nil {
		return err
	}

	ciphertext, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return json.Unmarshal(ciphertext, &v.Entries)
}
