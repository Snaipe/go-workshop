package vault

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
)

type Entry struct {
	Password string
}

type Vault struct {
	Entries map[string]Entry
}

func (v *Vault) Marshal(out io.Writer, block cipher.Block) error {
	plaintext, err := json.Marshal(v)
	if err != nil {
		return err
	}

	ciphertext := make([]byte, block.BlockSize()+len(plaintext))
	iv := ciphertext[:block.BlockSize()]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[block.BlockSize():], plaintext)

	_, err = out.Write(ciphertext)
	return err
}

func (v *Vault) Unmarshal(in io.Reader, block cipher.Block) error {
	ciphertext, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	if len(ciphertext) < block.BlockSize() {
		return fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:block.BlockSize()]
	ciphertext = ciphertext[block.BlockSize():]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	if err := json.Unmarshal(ciphertext, v); err != nil {
		return fmt.Errorf("could not decrypt vault; is the password correct?")
	}
	return nil
}
