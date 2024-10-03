package userdb

import (
	"context"

	"github.com/Zanda256/commitsmart-task/business/core/user"
	documentStore "github.com/Zanda256/commitsmart-task/business/data/docStore"
	"github.com/Zanda256/commitsmart-task/foundation/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

var dataKeyId = "base64 key"

// Store manages the set of APIs for rate search-db access.
type Store struct {
	log        *logger.Logger
	db         *mongo.Database
	userColl   *mongo.Collection
	encryptMgr *mongo.ClientEncryption
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *mongo.Database, encryptionClient *mongo.ClientEncryption, userCollectionName string) *Store {
	return &Store{
		log:        log,
		db:         db,
		userColl:   documentStore.OpenCollection(db, userCollectionName),
		encryptMgr: encryptionClient,
	}
}

// func (s *Store) encryptStrValue(val, dataKeyId string) error {
// 	nameRawValueType, nameRawValueData, err := bson.MarshalValue(val)
// 	if err != nil {
// 		panic(err)
// 	}
// 	nameRawValue := bson.RawValue{Type: nameRawValueType, Value: nameRawValueData}
// 	nameEncryptionOpts := options.Encrypt().
// 		SetAlgorithm("AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic").
// 		SetKeyID(dataKeyId)
// 	nameEncryptedField, err := s.encryptMgr.Encrypt(
// 		context.TODO(),
// 		nameRawValue,
// 		nameEncryptionOpts)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// Create inserts a new user into the database.
func (s *Store) Create(ctx context.Context, usr user.User) error {

	return nil
}
