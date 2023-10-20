package main

import (
	"context"
	"crypto/aes"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/term"
	"github.com/BurntSushi/toml"

	"example.com/paman/mongostore"
	"example.com/paman/vault"
)

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], fmt.Sprintf(format, args...))
	os.Exit(1)
}

func main() {
	var config struct {
		DefaultStore string `toml:"default_store"`

		Stores map[string]struct {
			Type string

			File struct {
				Path string
			}

			Mongo struct {
				URI string
			}
		}
	}

	var cli struct {
		Get      GetCmd      `cmd help:"retrieve a password"`
		Set      SetCmd      `cmd help:"assign a password"`
		Generate GenerateCmd `cmd help:"generate a password"`

		Store      string
		ConfigPath string `default:"$HOME/.config/paman.toml"`
	}
	ctx := kong.Parse(&cli)

	_, err := toml.DecodeFile(os.ExpandEnv(cli.ConfigPath), &config)
	switch {
	case errors.Is(err, os.ErrNotExist):
	case err != nil:
		fmt.Fprintf(os.Stderr, "warn: could not load config: %v", err)
	}

	fmt.Fprint(os.Stderr, "please enter your vault password: ")
	password, err := term.ReadPassword(0)
	if err != nil {
		fatalf("%v", err)
	}
	fmt.Fprintln(os.Stderr)

	salt := []byte(os.Getenv("USER"))

	const rounds = 1000

	key := pbkdf2.Key(password, salt, rounds, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		fatalf("%v", err)
	}

	var store vault.Store

	if cli.Store != "" {
		config.DefaultStore = cli.Store
	}

	switch cfg := config.Stores[config.DefaultStore]; cfg.Type {
	case "file":
		store = &vault.FileStore{
			Block: block,
			Path:  cfg.File.Path,
		}
	case "mongo":
		store, err = mongostore.NewStore(
			context.Background(),
			os.Getenv("USER"),
			block,
			options.Client().ApplyURI(cfg.Mongo.URI),
		)
		if err != nil {
			fatalf("%v", err)
		}
	}

	if store == nil {
		fatalf("store configuration not specified")
	}
	ctx.BindTo(store, (*vault.Store)(nil))

	if err := ctx.Run(); err != nil {
		fatalf("%v", err)
	}
}
