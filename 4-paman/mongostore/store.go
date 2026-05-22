package mongostore

import (
	"bytes"
	"context"
	"crypto/cipher"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"example.com/paman/vault"
)

type Store struct {
	client *mongo.Client
	block  cipher.Block
	email  string
}

type UserData struct {
	ID    primitive.ObjectID `bson:"_id"`
	Email string             `bson:"email"`
	Vault []byte             `bson:"vault"`
}

func NewStore(client *mongo.Client, block cipher.Block, email string) (*Store, error) {

	idx := client.Database("paman").Collection("vaults").Indexes()

	model := mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := idx.CreateOne(context.Background(), model); err != nil {
		return nil, err
	}

	return &Store{
		client: client,
		block:  block,
		email:  email,
	}, nil
}

func (s *Store) Load(v *vault.Vault) error {

	ctx := context.Background()
	coll := s.client.Database("paman").Collection("vaults")

	filter := bson.D{{"email", s.email}}

	var userdata UserData
	err := coll.FindOne(ctx, filter).Decode(&userdata)
	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		return os.ErrNotExist
	case err != nil:
		return err
	}

	return v.Unmarshal(bytes.NewReader(userdata.Vault), s.block)
}

func (s *Store) Store(v *vault.Vault) error {

	ctx := context.Background()

	var ciphertext bytes.Buffer
	if err := v.Marshal(&ciphertext, s.block); err != nil {
		return err
	}

	coll := s.client.Database("paman").Collection("vaults")

	filter := bson.D{{"email", s.email}}
	userdata := bson.D{{"$set", bson.D{{"vault", ciphertext.Bytes()}}}}

	_, err := coll.UpdateOne(ctx, filter, userdata, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

// NOTE: ceci permet de s'assurer que notre type *Store implémente
// l'interface vault.Store
var _ vault.Store = (*Store)(nil)
