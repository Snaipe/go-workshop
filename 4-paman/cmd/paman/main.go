package main

import (
	"context"
	"crypto/aes"
	"crypto/pbkdf2"
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/term"

	"example.com/paman/mongostore"
	"example.com/paman/vault"
)

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], fmt.Sprintf(format, args...))
	os.Exit(1)
}

func main() {
	var cli struct {
		Get      GetCmd      `cmd help:"retrieve a password"`
		Set      SetCmd      `cmd help:"assign a password"`
		Generate GenerateCmd `cmd help:"generate a password"`

		Password string `short:"p"`
	}
	ctx := kong.Parse(&cli)

	if cli.Password == "" {
		fmt.Fprint(os.Stderr, "please enter vault password: ")
		password, err := term.ReadPassword(0)
		fmt.Fprint(os.Stderr, "\n")
		if err != nil {
			fatalf("reading password: %v", err)
		}
		cli.Password = string(password)
	}

	// NOTE: dans la vraie vie, salt devrait être généré aléatoirement pour
	// chaque utilisateur, et stocké quelque part.
	salt := []byte(os.Getenv("USER"))

	// NOTE: le nombre de rounds de PBKDF2 choisi dépends du délai acceptable
	// de l'opération pour votre application. En pratique OWASP recommande un
	// minimum de 600 000 itérations pour HMAC_SHA_256.
	const rounds = 600_000

	key, err := pbkdf2.Key(sha256.New, cli.Password, salt, rounds, 32)
	if err != nil {
		fatalf("running key-derivation function: %v", err)
	}

	block, err := aes.NewCipher(key) // key défini plus tard
	if err != nil {
		fatalf("creating aes cipher: %v", err)
	}

	goctx := context.Background()
	uri := "mongodb://localhost:27017/"

	client, err := mongo.Connect(goctx, options.Client().ApplyURI(uri))
	if err != nil {
		fatalf("connecting to mongodb %q: %v", uri, err)
	}

	defer func() {
		// Attendre 5 minutes maximum que le client se déconnecte.
		ctx, stop := context.WithTimeout(context.Background(), 5*time.Minute)
		defer stop()

		if err := client.Disconnect(ctx); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	store, err := mongostore.NewStore(client, block, "user@example.com")
	if err != nil {
		fatalf("creating mongo store: %v", err)
	}

	//store := &vault.FileStore{
	//	Block: block, // sera défini plus tard
	//	Path:  "vault.dat",
	//}
	ctx.BindTo(store, (*vault.Store)(nil))

	if err := ctx.Run(); err != nil {
		fatalf("%v", err)
	}
}
