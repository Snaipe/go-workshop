package mongostore

import (
	"bytes"
	"context"
	"crypto/cipher"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"example.com/paman/vault"
)

type Store struct {
	block  cipher.Block
	ctx    context.Context
	client *mongo.Client
	user   string
}

func NewStore(ctx context.Context, user string, block cipher.Block, opts ...*options.ClientOptions) (*Store, error) {
	var (
		store Store
		err   error
	)
	store.block = block
	store.ctx = ctx
	store.user = user
	store.client, err = mongo.Connect(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (s *Store) Close() error {
	// Attendre 5 minutes maximum que le client se déconnecte.
	ctx, stop := context.WithTimeout(context.Background(), 5*time.Minute)
	defer stop()

	return s.client.Disconnect(ctx)
}

type vaultData struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	User  string
	Vault []byte
}

func (s *Store) Load(v *vault.Vault) error {
	coll := s.client.Database("paman").Collection("vaults")

	filter := bson.D{{"user", s.user}}
	var data vaultData
	err := coll.FindOne(s.ctx, filter).Decode(&data)

	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		*v = vault.Vault{
			Entries: map[string]vault.Entry{},
		}
		return nil
	case err != nil:
		return err
	}

	return v.Unmarshal(bytes.NewReader(data.Vault), s.block)
}

func (s *Store) Store(v *vault.Vault) error {
	var vdata bytes.Buffer
	if err := v.Marshal(&vdata, s.block); err != nil {
		return err
	}

	coll := s.client.Database("paman").Collection("vaults")
	data := vaultData{
		User:  s.user,
		Vault: vdata.Bytes(),
	}

	filter := bson.D{{"user", s.user}}
	opts := options.Replace().SetUpsert(true)
	_, err := coll.ReplaceOne(s.ctx, filter, data, opts)
	return err
}

// NOTE: ceci permet de s'assurer que notre type *Store implémente
// l'interface vault.Store
var _ vault.Store = (*Store)(nil)
