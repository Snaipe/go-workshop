package vault

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"maps"
	"testing"
)

func TestVaultMarshaling(t *testing.T) {
	cases := []struct {
		Name  string
		Vault Vault
	}{
		{
			Name:  "Empty",
			Vault: Vault{},
		},
		{
			Name: "Basic",
			Vault: Vault{
				Entries: map[string]Entry{
					"google": Entry{
						Password: "my google password",
					},
					"netflix": Entry{
						Password: "shared among friends",
					},
				},
			},
		},
	}

	testKey, _ := hex.DecodeString("6368616e676520746869732070617373")

	block, err := aes.NewCipher(testKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, tcase := range cases {
		t.Run(tcase.Name, func(t *testing.T) {
			var encrypted bytes.Buffer
			if err := tcase.Vault.Marshal(&encrypted, block); err != nil {
				t.Fatal(err)
			}
			var actual Vault
			if err := actual.Unmarshal(&encrypted, block); err != nil {
				t.Fatal(err)
			}
			if !maps.Equal(tcase.Vault.Entries, actual.Entries) {
				t.Fatalf("vaults are not equal: expected %v, got %v",
					tcase.Vault, actual)
			}
		})
	}
}
