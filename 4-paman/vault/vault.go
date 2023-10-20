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
	// Potentiellement d'autres champs plus tard
}

type Vault struct {
	Entries map[string]Entry
}

func (v *Vault) Marshal(out io.Writer, block cipher.Block) error {
	// FIXME:
	// 1. utilisez encoding/json pour transformer Vault en texte à chiffrer
	// 2. utilisez block pour chiffrer le texte
	// 3. écrivez le texte dans out
	//
	// NOTE: l'utilisation de cipher.Block n'est pas triviale; ne tentez pas
	// de le deviner par vous même. Référez vous à l'exemple dans la bibliothèque
	// standard: https://pkg.go.dev/crypto/cipher#NewCFBEncrypter

	plaintext, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
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
	// FIXME:
	// 1. lisez le document chiffré depuis le Reader
	// 2. utilisez block pour déchiffrer le document
	// 3. utilisez encoding/json pour décoder le texte en Vault
	//
	// NOTE: l'utilisation de cipher.Block n'est pas triviale; ne tentez pas
	// de le deviner par vous même. Référez vous à l'exemple dans la bibliothèque
	// standard: https://pkg.go.dev/crypto/cipher#NewCFBDecrypter

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
